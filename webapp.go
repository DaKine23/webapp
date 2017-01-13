package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "expvar"

	"github.com/DaKine23/webapp/hb"
	"github.com/DaKine23/webapp/hb/bsbutton"
	"github.com/DaKine23/webapp/hb/bsbuttongroup"
	"github.com/DaKine23/webapp/hb/bsglyphicons"

	"github.com/DaKine23/webapp/hb/bstable"
	"github.com/DaKine23/webapp/hb/jqaction"
	"github.com/DaKine23/webapp/hb/svg"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	cors "github.com/itsjamie/gin-cors"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

func main() {
	flag.Parse()
	initData()

	configRuntime()
	go configMonitoring()

	configController()
}

var data = make(map[string]*[]*hb.HTMLTableRow)
var titles = make(map[string][]string)
var pagesize = make(map[string]int)

func initData() {
	titles["mytable"] = []string{"ping or pong", "timestamp", "delete", "id"}
	data["mytable"] = &[]*hb.HTMLTableRow{}
	pagesize["mytable"] = 10
	titles["mytable2"] = []string{"one", "two", "three"}
	pagesize["mytable2"] = 10
	data["mytable2"] = &[]*hb.HTMLTableRow{
		&hb.HTMLTableRow{
			Row: &[]interface{}{"asome", "rontent", 1},
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"bknow", "yblubb", 3},
			Status: bstable.TableRowStatusDanger,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"csome", "xcontent", 2},
			Status: bstable.TableRowStatusInfo,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"dknow", "cblubb", 6},
			Status: bstable.TableRowStatusSuccess,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"esome", "fcontent", 5},
			Status: bstable.TableRowStatusWarning,
		},
		&hb.HTMLTableRow{
			Row: &[]interface{}{"fknow", "ablubb", 42},
		},
		&hb.HTMLTableRow{
			Row: &[]interface{}{"asome", "rontent", 1},
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"bknow", "yblubb", 3},
			Status: bstable.TableRowStatusDanger,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"csome", "xcontent", 2},
			Status: bstable.TableRowStatusInfo,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"dknow", "cblubb", 6},
			Status: bstable.TableRowStatusSuccess,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"esome", "fcontent", 5},
			Status: bstable.TableRowStatusWarning,
		},
		&hb.HTMLTableRow{
			Row: &[]interface{}{"fknow", "ablubb", 42},
		},
	}
}

// responseJSON provides a stucture for a JSON message response
type responseJSON struct {
	Message string    `json:"message" xml:"message"`
	Time    time.Time `json:"time" xml:"time"`
}

type responseTableJSON struct {
	Table string `json:"table" xml:"table"`
}

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	glog.Infof("Running with %d CPUs\n", nuCPU)
}

var switchingValue string

var id int

func switchValue() {

	if switchingValue == "ping" {
		switchingValue = "pong"
	} else {
		switchingValue = "ping"
	}

}

func everyXSecondsDo(sec int, dof do) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	for {
		select {
		case <-ticker.C:
			dof()
		}
	}
}

type do func()

func configMonitoring() {
	http.ListenAndServe(":1234", nil)
}

func configController() {
	router := gin.Default()
	// router.Use(ginglog.Logger(3 * time.Second))
	// router.Use(ginoauth2.RequestLogger([]string{"uid"}, "data"))
	// router.Use(gin.Recovery())

	ginoauth2.VarianceTimer = 300 * time.Millisecond // defaults to 30s

	//router.LoadHTMLGlob("public/*.html")
	router.Static("/static", "./public/static")
	router.Static("/fonts", "./public/fonts")

	go everyXSecondsDo(1, switchValue)

	//add cors
	addCors(router)

	handler := func(c *gin.Context) {
		//c.String(http.StatusOK,"Hello World")

		c.JSON(http.StatusOK, responseJSON{switchingValue, time.Now()})
	}

	router.GET("/pong", handler)
	//GET http://localhost:3000

	nameParamHandler := func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	}

	router.GET("/hello/:name", nameParamHandler)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "index.html")
	})
	router.RedirectTrailingSlash = true

	router.GET("/index.html", func(c *gin.Context) {
		c.Writer.WriteString(page())

		//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
		//c.HTML(200, "index.html", gin.H{})
	})

	router.POST("/table/:tablename/add/:count", addNewLineToUpperTableHandler)
	router.POST("/table/:tablename/", addNewLineToLowerTableHandler)

	router.DELETE("/table/:tablename/delete/:id/:page", deleteHandler)

	router.GET("/table/:tablename/sort/:column/:page", sortHandler)
	router.GET("/table/:tablename/show/:page", showHandler)

	router.StaticFile("/favicon.ico", "public/favicon.ico")

	router.Run(":3000")
}

func addCors(router *gin.Engine) {
	//if you want to use cors
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
}

func page() string {

	//define doctype
	result := "<!DOCTYPE html svg>\n"

	//define <html>
	html := hb.NewHTMLPart("html", "", "")

	//define <head> with title and include bootstrap
	title := "Webapp Example"
	head := hb.NewHTMLPart("head", "", `<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="static/bootstrap.min.css">
	<link rel="stylesheet" href="static/font-awesome.min.css">
	<link rel="stylesheet" href="static/zalos-bootstrap-theme.min.css">`).
		AddSubParts(hb.NewHTMLPart("title", "", title))

	head.AddSubParts(hb.NewHTMLPart("style", "", `.icon {
  width: 14px;
  height: 14px;
}
	`))
	//define js libraries you want to import
	jsLibraries := []string{

		"static/jquery.min.js",
		"static/bootstrap.min.js",
	}
	//add js libraries you want to import to the <head>
	for _, v := range jsLibraries {
		jsLibrariesPart := hb.NewHTMLPart("script", "", "").AddOption(&hb.HTMLOption{
			Name:  "src",
			Value: v,
		})
		head.AddSubParts(jsLibrariesPart)
	}

	//define <body>
	body := hb.NewHTMLPart("body", "", svg.Iconset)

	// define two tables for demonstration
	table := hb.NewHTMLTable("mytable", titles["mytable"], *data["mytable"], []string{}, pagesize["mytable"], 1)
	table2 := hb.NewHTMLTable("mytable2", titles["mytable2"], *data["mytable2"], []string{}, pagesize["mytable2"], 1)

	// tables should be used inside of containers when defining the layout will be used as drawing destination later on
	tp := hb.NewHTMLTableContainer(table)
	tp2 := hb.NewHTMLTableContainer(table2)

	// some input fields for table 2 to enter new rows
	edit := hb.NewLineEdit("myinput", "Mighty Input", "may type sth here", "standard content", nil)
	edit2 := hb.NewLineEdit("myinput2", hb.NewGlyphicon(bsglyphicons.GlyphiconEurGlyphiconEuro).String(), "money money money", "", nil)
	searchedit := hb.NewLineEdit("myinput3", hb.NewGlyphicon(bsglyphicons.GlyphiconBook).String(), "Search some thing", "", nil).AddLineEditSearchButton("searchbutton")
	submitbutton := hb.NewHTMLPart("button", "submitbutton", "submit").AddBootstrapClasses(bsbutton.B, bsbutton.Primary, bsbutton.BlockLevel)

	// add some buttons  ("a" for bootstrap buttongroups)
	button := hb.NewHTMLPart("a", "addbutton", "Add Table Entry").AddBootstrapClasses(bsbutton.B, bsbutton.Primary)
	button3 := hb.NewHTMLPart("a", "addbutton2", "Add 1000 Table Entries").AddBootstrapClasses(bsbutton.B, bsbutton.Primary)
	button2 := hb.NewHTMLPart("button", "pongbutton", "Click Me").AddBootstrapClasses(bsbutton.B, bsbutton.Default)

	//Create a buttongroup
	buttongroup := hb.NewHTMLPart("div", "", "").
		AddBootstrapClasses(bsbuttongroup.ButtonGroup, bsbuttongroup.JustifiedButtonGroup).
		AddSubParts(button, button3)

	// add some html <div>
	div2 := hb.NewHTMLPart("div", "drawdestination", "content")

	// define scripts (ajax calls)
	scriptToAddARow := hb.NewScript(button.ID, jqaction.Click, table.ID+"container", "POST", "/table/mytable/add/1", hb.JSONResultValue("table"))
	scriptToAdd1000Rows := hb.NewScript(button3.ID, jqaction.Click, table.ID+"container", "POST", "/table/mytable/add/1000", hb.JSONResultValue("table"))
	script := hb.NewScript(button2.ID, jqaction.Click, div2.ID, "GET", "/pong", hb.JSONResultValue("message")+`+" !!! " +`+hb.JSONResultValue("timestamp"))

	// add <head> and <body> to <html>
	html.AddSubParts(head, body)

	//create a bootstrap grid

	numberedit := hb.NewLineEdit("myinput4", "Numbers", "12.4", "", nil)

	validator := "^[1-9]+[0-9]*$"
	intgeredit2 := hb.NewLineEdit("myinput5", hb.NewGlyphicon(bsglyphicons.GlyphiconKnight).String(), "12345", "", &hb.Validation{RegEx: validator})

	intgeredit2.AddTooltip("I only accept Integer Values", "left")
	body.AddScripts(hb.TooltipScript())

	somediv := hb.NewHTMLPart("div", "keypressdiv", "here")
	keypressscript := hb.NewHTMLPart("script", "", `$(document).ready(function(){$("#myinput4").keypress(function(event){
    $("#keypressdiv").html("Key: " + event.which);
	});});`).AddOption(&hb.HTMLOption{Name: "type", Value: "text/javascript"})

	var onerr string
	onsuc := hb.OnResult(tp2.ID, hb.JSONResultValue("table")) + `$("#myinput").focus();$("#myinput").val('');$("#myinput2").val('');$("#myinput3").val('');`
	ig := hb.InputGroup{
		Member: []hb.InputGroupMember{{"myinput", "one"}, {"myinput2", "two"}, {"myinput3", "three"}, {"dropdownsample", "status"}},
	}

	igscript := hb.NewInputGroupScript(submitbutton.ID, jqaction.Click, "", "POST", "/table/"+table2.ID, ig, onsuc, onerr)
	igscript2 := hb.NewInputGroupScript("myinput3", jqaction.Keypress, "event.which == 13", "POST", "/table/"+table2.ID, ig, onsuc, onerr)

	// add all the other html tags to the <body>

	dropdown := hb.NewDropDownInput("dropdownsample", hb.NewGlyphicon(bsglyphicons.GlyphiconEducation).String(), false, "active", "success", "info", "warning", "danger")

	grid := hb.BsGrid{&[][]hb.BsCell{
		{{buttongroup, 0}},
		{{nil, 2}, {tp, 8}},
		{{button2, 0}},
		{{div2, 0}},
		{{edit, 0}, {edit2, 0}, {searchedit, 0}, {dropdown, 0}, {submitbutton, 0}},
		{{nil, 2}, {tp2, 8}},
		{{numberedit, 0}, {intgeredit2, 0}},
	}, ""}

	body.AddSubParts(script.HTMLPart, scriptToAddARow.HTMLPart, scriptToAdd1000Rows.HTMLPart, grid.HTMLPart(), somediv)
	body.AddScripts(keypressscript, igscript, igscript2)

	// return DOCTYPE definition + <html> as string (includes all the subparts)
	return result + html.String()

}

// hb tables enforce "/table/:tablename/sort/:column" api
// hb.Sort helps you sort as it sorts the rows using the value on provided index
// and discovers integer float64 and time.Time or uses the string representation for sorting
func sortHandler(c *gin.Context) {

	tbn := c.Param("tablename")
	tbc := c.Param("column")
	var tbp int
	if itbp, err := strconv.Atoi(c.Param("page")); err == nil {
		tbp = itbp
	}

	//go through titles and find the corresponding column
	for i, v := range titles[tbn] {
		//shorten the columns title by convention
		reducedTitle := strings.Replace(fmt.Sprint(v), " ", "", -1)
		reducedTitle = strings.ToLower(reducedTitle)

		// if found sort the data
		if tbc == reducedTitle {

			*data[tbn] = hb.Sort(i, *data[tbn])

			break
		}

	}

	//draw the table
	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{}, pagesize[tbn], tbp)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func showHandler(c *gin.Context) {
	tbn := c.Param("tablename")
	var tbp int
	if itbp, err := strconv.Atoi(c.Param("page")); err == nil {
		tbp = itbp
	}
	//draw the table
	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{}, pagesize[tbn], tbp)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func deleteHandler(c *gin.Context) {

	id := c.Param("id")
	tbn := c.Param("tablename")
	var tbp int
	if itbp, err := strconv.Atoi(c.Param("page")); err == nil {
		tbp = itbp
	}

	index := 0

	//get the index of the "id column"
	for i, v := range *(*data[tbn])[0].ParentTable.Titles.Row {
		if v == "id" {
			index = i
			break
		}

	}
	//remove row search on row index of the "id column"
	for i, v := range *data[tbn] {
		if fmt.Sprint((*v.Row)[index]) == id {
			*data[tbn] = append((*data[tbn])[:i], (*data[tbn])[i+1:]...)
			break
		}
	}

	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{}, pagesize[tbn], tbp)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func addNewLineToLowerTableHandler(c *gin.Context) {
	//read table name from uri
	tbn := c.Param("tablename")

	type Data struct {
		//exists checks if field is part of the request
		One string `json:"one" binding:"exists"`
		//required checks if field is part of the request and not empty
		Two    string `json:"two" binding:"required"`
		Three  string `json:"three" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	lastpage := len(*data[tbn]) / pagesize[tbn]
	if len(*data[tbn])%pagesize[tbn] != 0 {
		lastpage++
	}

	var json Data
	if c.BindJSON(&json) == nil {
		newrow := []interface{}{json.One, json.Two, json.Three}
		*data[tbn] = append(*data[tbn], hb.NewHTMLTableRow(newrow, json.Status))
		table2 := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{}, pagesize[tbn], lastpage)

		c.JSON(http.StatusOK, gin.H{"table": table2.String()})

	}

}

func addNewLineToUpperTableHandler(c *gin.Context) {

	//read table name from uri
	tbn := c.Param("tablename")
	count, _ := strconv.Atoi(c.Param("count"))

	for i := 0; i < count; i++ {
		//increment serial (would be done on the database usually)
		id++

		ids := fmt.Sprint(id)
		//create a HTMLPart to hold a button and script
		buttoncontainer := hb.NewHTMLPart("deletebutton", "", "")

		// create a Bootstrap styled Button
		button := hb.NewHTMLPart("button", "tablebutton"+ids, "del").AddBootstrapClasses(bsbutton.B, bsbutton.SizeVerySmall, bsbutton.Danger)
		// create a deletion script for the Button to delete the row containing the button
		script := hb.NewTableButtonScript(button.ID, jqaction.Click, tbn+"container", tbn, "DELETE", "table/"+tbn+"/delete/"+ids, hb.JSONResultValue("table"))
		// add button and script to the container
		buttoncontainer.AddSubParts(script.HTMLPart, button)

		// data creation and appending could may be done on the database
		// create a new row to insert
		newrow := []interface{}{switchingValue, time.Now(), buttoncontainer.String(), id}
		// append it to the data
		*data[tbn] = append(*data[tbn], hb.NewHTMLTableRow(newrow))
	}
	//create table using the Data and return the HTML

	lastpage := len(*data[tbn]) / pagesize[tbn]
	if len(*data[tbn])%pagesize[tbn] != 0 {
		lastpage++
	}

	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{}, pagesize[tbn], lastpage)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}
