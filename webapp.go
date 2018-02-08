package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "expvar"

	"github.bus.zalan.do/ale/gowt"
	"github.bus.zalan.do/ale/gowt/bsbutton"
	"github.bus.zalan.do/ale/gowt/bsglyphicons"

	"github.bus.zalan.do/ale/gowt/bstable"
	"github.bus.zalan.do/ale/gowt/jqaction"
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

var data = make(map[string][]gowt.HTMLTableRow)
var titles = make(map[string][]string)
var pagesize = make(map[string]int)

func initData() {
	titles["mytable"] = []string{"ping or pong", "timestamp", "delete", "id"}
	data["mytable"] = []gowt.HTMLTableRow{}
	pagesize["mytable"] = 10
	titles["mytable2"] = []string{"one", "two", "three"}
	pagesize["mytable2"] = 10
	data["mytable2"] = []gowt.HTMLTableRow{
		gowt.HTMLTableRow{
			Row: []interface{}{"asome", "rontent", 1},
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"bknow", "yblubb", 3},
			Status: bstable.TableRowStatusDanger,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"csome", "xcontent", 2},
			Status: bstable.TableRowStatusInfo,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"dknow", "cblubb", 6},
			Status: bstable.TableRowStatusSuccess,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"esome", "fcontent", 5},
			Status: bstable.TableRowStatusWarning,
		},
		gowt.HTMLTableRow{
			Row: []interface{}{"fknow", "ablubb", 42},
		},
		gowt.HTMLTableRow{
			Row: []interface{}{"asome", "rontent", 1},
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"bknow", "yblubb", 3},
			Status: bstable.TableRowStatusDanger,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"csome", "xcontent", 2},
			Status: bstable.TableRowStatusInfo,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"dknow", "cblubb", 6},
			Status: bstable.TableRowStatusSuccess,
		},
		gowt.HTMLTableRow{
			Row:    []interface{}{"esome", "fcontent", 5},
			Status: bstable.TableRowStatusWarning,
		},
		gowt.HTMLTableRow{
			Row: []interface{}{"fknow", "ablubb", 42},
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
	})

	router.POST("/table/:tablename/add/:count", addNewLineToUpperTableHandler)
	router.POST("/table/:tablename/", addNewLineToLowerTableHandler)

	router.DELETE("/table/:tablename/delete/:id/:page", deleteHandler)

	router.GET("/table/:tablename/sort/:column/:page", sortHandler)
	router.GET("/table/:tablename/show/:page", showHandler)
	router.GET("/queries", queryhandler)

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

	html, _, body := gowt.NewDefaultPage("Webapp example", []string{}, []string{"static/ant-strap.min.css"})

	// define two tables for demonstration
	table := gowt.NewHTMLTable("mytable", titles["mytable"], data["mytable"], []string{}, pagesize["mytable"], 1)
	table2 := gowt.NewHTMLTable("mytable2", titles["mytable2"], data["mytable2"], []string{}, pagesize["mytable2"], 1)

	// tables should be used inside of containers when defining the layout will be used as drawing destination later on
	tp := gowt.NewHTMLTableContainer(table)
	tp2 := gowt.NewHTMLTableContainer(table2)

	// some input fields for table 2 to enter new rows
	edit := gowt.NewLineEdit("myinput", "Mighty Input", "may type sth here", "standard content", nil)
	edit2 := gowt.NewLineEdit("myinput2", gowt.NewGlyphicon(bsglyphicons.GlyphiconEurGlyphiconEuro).String(), "money money money", "", nil)
	searchedit := gowt.NewLineEdit("myinput3", gowt.NewGlyphicon(bsglyphicons.GlyphiconBook).String(), "Search some thing", "", nil).AddLineEditSearchButton("searchbutton")
	submitbutton := gowt.NewHTMLPart("button", "submitbutton", "submit").AddBootstrapClasses(bsbutton.B, bsbutton.Primary, bsbutton.BlockLevel)

	// add some buttons  ("a" for bootstrap buttongroups)
	button := gowt.NewHTMLPart("a", "addbutton", "Add Table Entry").AddOption(&gowt.HTMLOption{"role", "button"}).AddBootstrapClasses(bsbutton.B, bsbutton.Primary, bsbutton.BlockLevel)
	button3 := gowt.NewHTMLPart("a", "addbutton2", "Add 1000 Table Entries").AddOption(&gowt.HTMLOption{"role", "button"}).AddBootstrapClasses(bsbutton.B, bsbutton.Primary, bsbutton.BlockLevel)
	button2 := gowt.NewHTMLPart("button", "pongbutton", "Click Me").AddBootstrapClasses(bsbutton.B, bsbutton.Default)

	// add some html <div>
	div2 := gowt.NewHTMLPart("div", "drawdestination", "content")

	// define scripts (ajax calls)
	scriptToAddARow := gowt.NewScript(button.ID, jqaction.Click, table.ID+"container", "POST", "/table/mytable/add/1", gowt.JSONResultValue("table"))
	scriptToAdd1000Rows := gowt.NewScript(button3.ID, jqaction.Click, table.ID+"container", "POST", "/table/mytable/add/1000", gowt.JSONResultValue("table"))
	script := gowt.NewScript(button2.ID, jqaction.Click, div2.ID, "GET", "/pong", gowt.JSONResultValue("message")+`+" !!! " +`+gowt.JSONResultValue("timestamp"))

	//create a bootstrap grid

	numberedit := gowt.NewLineEdit("myinput4", "Numbers", "12.4", "", nil)

	validator := "^[1-9]+[0-9]*$"
	intgeredit2 := gowt.NewLineEdit("myinput5", gowt.NewGlyphicon(bsglyphicons.GlyphiconKnight).String(), "12345", "", &gowt.Validation{RegEx: validator})

	intgeredit2.AddTooltip("I only accept Integer Values", "left")
	body.AddScripts(gowt.TooltipScript())

	somediv := gowt.NewHTMLPart("div", "keypressdiv", "here")

	// js for displaying the current keys number
	keypressscript := gowt.NewHTMLPart("script", "", `$(document).ready(function(){$("#myinput4").keypress(function(event){
    $("#keypressdiv").html("Key: " + event.which);
	});});`).AddOption(&gowt.HTMLOption{Name: "type", Value: "text/javascript"})

	var onerr string

	// js for set focus to first input and clean values
	onsuc := gowt.OnResult(tp2.ID, gowt.JSONResultValue("table")) + `$("#myinput").focus();$("#myinput").val('');$("#myinput2").val('');$("#myinput3").val('');`
	ig := gowt.InputGroup{
		Member: []gowt.InputGroupMember{{"myinput", "one"}, {"myinput2", "two"}, {"myinput3", "three"}, {"dropdownsample", "status"}},
	}

	// InputGroup Scripts are jquery that is called when defined source is triggered by defined action. it sends a request with defined resttype containing the values in the defined inputs as json with defined valuenames
	// onSuccess or onError are executed in addition! can be used to warn or inform the user clean inputs and so on
	igscript := gowt.NewInputGroupScript(submitbutton.ID, jqaction.Click, "", "GET", "/table/"+table2.ID, ig, onsuc, onerr)
	igscript2 := gowt.NewInputGroupScript("myinput3", jqaction.Keypress, "event.which == 13", "POST", "/table/"+table2.ID, ig, onsuc, onerr)

	// add all the other html tags to the <body>

	dropdown := gowt.NewDropDownInput("dropdownsample", gowt.NewGlyphicon(bsglyphicons.GlyphiconEducation).String(), false, "active", "success", "info", "warning", "danger")

	grid := gowt.BsGrid{&[][]gowt.BsCell{
		{{button, 0}, {button3, 0}},
		{{nil, 2}, {tp, 8}},
		{{button2, 0}},
		{{div2, 0}},
		{{edit, 0}, {edit2, 0}, {searchedit, 0}, {dropdown, 0}, {submitbutton, 0}},
		{{nil, 2}, {tp2, 8}},
		{{numberedit, 0}, {intgeredit2, 0}},
	}, ""}

	body.AddSubParts(grid.HTMLPart(), somediv)
	body.AddScripts(keypressscript, igscript, igscript2, scriptToAddARow, script, scriptToAdd1000Rows)

	// return DOCTYPE definition + <html> as string (includes all the subparts)
	return result + html.String()

}

// hb tables enforce "/table/:tablename/sort/:column" api
// gowt.Sort helps you sort as it sorts the rows using the value on provided index
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

			data[tbn] = gowt.Sort(i, data[tbn])

			break
		}

	}

	//draw the table
	table := gowt.NewHTMLTable(tbn, titles[tbn], data[tbn], []string{}, pagesize[tbn], tbp)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func queryhandler(c *gin.Context) {
	req := c.Request
	unesc, _ := url.QueryUnescape(req.URL.RawQuery)

	c.String(200, unesc)
}

func showHandler(c *gin.Context) {
	tbn := c.Param("tablename")
	var tbp int
	if itbp, err := strconv.Atoi(c.Param("page")); err == nil {
		tbp = itbp
	}
	//draw the table
	table := gowt.NewHTMLTable(tbn, titles[tbn], data[tbn], []string{}, pagesize[tbn], tbp)
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
	for i, v := range data[tbn][0].ParentTable.Titles.Row {
		if v == "id" {
			index = i
			break
		}

	}
	//remove row search on row index of the "id column"
	for i, v := range data[tbn] {
		if fmt.Sprint((v.Row)[index]) == id {
			data[tbn] = append(data[tbn][:i], data[tbn][i+1:]...)
			break
		}
	}

	table := gowt.NewHTMLTable(tbn, titles[tbn], data[tbn], []string{}, pagesize[tbn], tbp)
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

	lastpage := len(data[tbn]) / pagesize[tbn]
	if len(data[tbn])%pagesize[tbn] != 0 {
		lastpage++
	}

	var json Data
	if c.BindJSON(&json) == nil {
		newrow := []interface{}{json.One, json.Two, json.Three}
		data[tbn] = append(data[tbn], *gowt.NewHTMLTableRow(newrow, json.Status))
		table2 := gowt.NewHTMLTable(tbn, titles[tbn], data[tbn], []string{}, pagesize[tbn], lastpage)

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
		buttoncontainer := gowt.NewHTMLPart("deletebutton", "")

		// create a Bootstrap styled Button
		button := gowt.NewHTMLPart("button", "tablebutton"+ids, "del").AddBootstrapClasses(bsbutton.B, bsbutton.SizeVerySmall, bsbutton.Danger)
		// create a deletion script for the Button to delete the row containing the button
		script := gowt.NewTableButtonScript(button.ID, jqaction.Click, tbn+"container", tbn, "DELETE", "table/"+tbn+"/delete/"+ids, gowt.JSONResultValue("table"))
		// add button and script to the container
		buttoncontainer.AddSubParts(button)
		buttoncontainer.AddScripts(script)

		// data creation and appending could may be done on the database
		// create a new row to insert
		newrow := []interface{}{switchingValue, time.Now(), buttoncontainer.String(), id}
		// append it to the data
		data[tbn] = append(data[tbn], *gowt.NewHTMLTableRow(newrow))
	}
	//create table using the Data and return the HTML

	lastpage := len(data[tbn]) / pagesize[tbn]
	if len(data[tbn])%pagesize[tbn] != 0 {
		lastpage++
	}

	table := gowt.NewHTMLTable(tbn, titles[tbn], data[tbn], []string{}, pagesize[tbn], lastpage)
	asc := true
	table.ToSortableTable("timestamp", &asc)
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}
