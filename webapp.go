package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "expvar"
	"flag"

	"github.com/DaKine23/webapp/hb"
	"github.com/DaKine23/webapp/hb/bsbutton"
	"github.com/DaKine23/webapp/hb/bsbuttongroup"
	"github.com/DaKine23/webapp/hb/bscontainer"
	"github.com/DaKine23/webapp/hb/bsglyphicons"
	"github.com/DaKine23/webapp/hb/bsgrid"
	"github.com/DaKine23/webapp/hb/bstable"
	"github.com/DaKine23/webapp/hb/jqaction"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	cors "github.com/itsjamie/gin-cors"
	ginglog "github.com/zalando/gin-glog"
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
	router := gin.New()
	router.Use(ginglog.Logger(3 * time.Second))
	router.Use(ginoauth2.RequestLogger([]string{"uid"}, "data"))
	router.Use(gin.Recovery())

	ginoauth2.VarianceTimer = 300 * time.Millisecond // defaults to 30s

	//router.LoadHTMLGlob("public/*.html")
	//router.Static("/static", "./public/static")

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
	result := "<!DOCTYPE html>\n"

	//define <html>
	html := hb.NewHTMLPart("html", "", "")

	//define <head> with title and include bootstrap
	title := "Webapp Example"
	head := hb.NewHTMLPart("head", "", `<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">`).
		AddSubParts(hb.NewHTMLPart("title", "", title))

		//define js libraries you want to import
	jsLibraries := []string{

		"https://unpkg.com/babel-standalone@6.15.0/babel.min.js",
		"https://unpkg.com/jquery@3.1.0/dist/jquery.min.js",
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
	body := hb.NewHTMLPart("body", "", "")

	// define two tables for demonstration
	table := hb.NewHTMLTable("mytable", titles["mytable"], *data["mytable"], []string{}, pagesize["mytable"], 1)
	table2 := hb.NewHTMLTable("mytable2", titles["mytable2"], *data["mytable2"], []string{}, pagesize["mytable2"], 1)

	// tables should be used inside of containers when defining the layout will be used as drawing destination later on
	tp := hb.NewHTMLTableContainer(table)
	tp2 := hb.NewHTMLTableContainer(table2)

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
	root := hb.NewHTMLPart("root", "", "").AddBootstrapClasses(bscontainer.ContainerFluid)
	row1 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	cell11 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	cell11.AddSubParts(buttongroup)
	row1.AddSubParts(cell11)

	row2 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	cell21 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	cell21.AddSubParts(tp)
	row2.AddSubParts(cell21)

	row3 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	cell31 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	cell31.AddSubParts(button2)
	row3.AddSubParts(cell31)

	row4 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	cell41 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	cell41.AddSubParts(div2)
	row4.AddSubParts(cell41)

	row5 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	cell51 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	cell51.AddSubParts(tp2)
	row5.AddSubParts(cell51)

	row6 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)

	cell61 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(4, bsgrid.Large))
	cell62 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(4, bsgrid.Large))
	cell63 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(4, bsgrid.Large))

	edit := hb.NewLineEdit("myinput", "Mighty Input", "may type sth here", "standard content")
	edit2 := hb.NewLineEdit("myinput2", hb.NewGlyphicon(bsglyphicons.GlyphiconEurGlyphiconEuro).String(), "money money money", "")
	searchedit := hb.NewLineEdit("myinput2", "", "Search some thing", "").AddLineEditSearch("searchbutton")

	cell61.AddSubParts(edit)
	cell62.AddSubParts(edit2)
	cell63.AddSubParts(searchedit)
	row6.AddSubParts(cell61, cell62, cell63)

	root.AddSubParts(row1, row2, row3, row4, row5, row6)

	// add all the other html tags to the <body>
	body.AddSubParts(script.HTMLPart, scriptToAddARow.HTMLPart, scriptToAdd1000Rows.HTMLPart, root)

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
