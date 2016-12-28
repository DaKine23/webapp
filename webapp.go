package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	_ "expvar"
	"flag"

	"github.com/DaKine23/webapp/hb"
	bsbutton "github.com/DaKine23/webapp/hb/bsbutton"
	bstable "github.com/DaKine23/webapp/hb/bstable"
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

func initData() {
	titles["mytable"] = []string{"ping or pong", "timestamp", "delete", "id"}
	data["mytable"] = &[]*hb.HTMLTableRow{}
	titles["mytable2"] = []string{"one", "two", "three"}
	//sample data
	data["mytable2"] = &[]*hb.HTMLTableRow{
		&hb.HTMLTableRow{
			Row: &[]interface{}{"asome", "rontent", 1},
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"bknow", "yblubb", 3},
			Status: bstable.BsTableRowStatusDanger,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"csome", "xcontent", 2},
			Status: bstable.BsTableRowStatusInfo,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"dknow", "cblubb", 6},
			Status: bstable.BsTableRowStatusSuccess,
		},
		&hb.HTMLTableRow{
			Row:    &[]interface{}{"esome", "fcontent", 5},
			Status: bstable.BsTableRowStatusWarning,
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

//the result type JSON schema the script expects
var tableResultType = hb.JSONResult{
	Names: []hb.JSONResultName{{"table"}},
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

	router.POST("/table/:tablename/add", addNewLineToUpperTableHandler)

	router.DELETE("/table/:tablename/delete/:id", deleteHandler)

	router.GET("/table/:tablename/sort/:column", sortHandler)

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

	//define json schema a script expects when called
	pongResultType := hb.JSONResult{
		Names: []hb.JSONResultName{{"message"}, {"time"}}, //just message and time for demonstration
	}
	//define json schema a script expects when called
	tableResultType := hb.JSONResult{
		Names: []hb.JSONResultName{{"table"}}, //contains table as html string
	}

	// define two tables for demonstration
	table := hb.NewHTMLTable("mytable", titles["mytable"], *data["mytable"], []string{})
	table2 := hb.NewHTMLTable("mytable2", titles["mytable2"], *data["mytable2"], []string{})

	// tables should be used inside of containers when defining the layout will be used as drawing destination later on
	tp := hb.NewHTMLTableContainer(table)
	tp2 := hb.NewHTMLTableContainer(table2)

	// add some buttons
	button := hb.NewHTMLPart("button", "addbutton", "Add Table Entry").AddBootstrapClasses(bsbutton.BsButton, bsbutton.BsButtonPrimary)
	button2 := hb.NewHTMLPart("button", "pongbutton", "Click Me").AddBootstrapClasses(bsbutton.BsButton, bsbutton.BsButtonDefault)

	// add some html <div>
	div2 := hb.NewHTMLPart("div", "drawdestination", "content")

	// define scripts (ajax calls)
	scriptToAddARow := hb.NewScript(button.ID, "click", table.ID+"container", "POST", "/table/mytable/add", tableResultType, tableResultType.Names[0].Value())
	script := hb.NewScript(button2.ID, "click", div2.ID, "GET", "/pong", pongResultType, pongResultType.Names[0].Value()+`+" !!! " +`+pongResultType.Names[1].Value())

	// add <head> and <body> to <html>
	html.AddSubParts(head, body)

	// add all the other html tags to the <body>
	body.AddSubParts(script.HTMLPart, scriptToAddARow.HTMLPart, button, tp, button2, div2, tp2)

	// return DOCTYPE definition + <html> as string (includes all the subparts)
	return result + html.String()

}

// hb tables enforce "/table/:tablename/sort/:column" api
// hb.Sort helps you sort as it sorts the rows using the value on provided index
// and discovers integer float64 and time.Time or uses the string representation for sorting
func sortHandler(c *gin.Context) {

	tbn := c.Param("tablename")
	tbc := c.Param("column")

	for i, v := range titles[tbn] {
		reducedTitle := strings.Replace(fmt.Sprint(v), " ", "", -1)
		reducedTitle = strings.ToLower(reducedTitle)

		if tbc == reducedTitle {

			*data[tbn] = hb.Sort(i, *data[tbn])

			break
		}

	}

	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{})
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func deleteHandler(c *gin.Context) {

	id := c.Param("id")
	tbn := c.Param("tablename")
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

	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{})
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}

func addNewLineToUpperTableHandler(c *gin.Context) {

	//read table name from uri
	tbn := c.Param("tablename")
	//for i := 0; i < 100; i++ {
	//increment serial (would be done on the database usually)
	id++

	ids := fmt.Sprint(id)
	//create a HTMLPart to hold a button and script
	buttoncontainer := hb.NewHTMLPart("deletebutton", "", "")

	// create a Bootstrap styled Button
	button := hb.NewHTMLPart("button", "tablebutton"+ids, "del").AddBootstrapClasses(bsbutton.BsButton, bsbutton.BsButtonDanger)
	// create a deletion script for the Button to delete the row containing the button
	script := hb.NewScript(button.ID, "click", tbn+"container", "DELETE", "table/"+tbn+"/delete/"+ids, tableResultType, tableResultType.Names[0].Value())
	// add button and script to the container
	buttoncontainer.AddSubParts(script.HTMLPart, button)

	// data creation and appending could may be done on the database
	// create a new row to insert
	newrow := []interface{}{switchingValue, time.Now(), buttoncontainer.String(), id}
	// append it to the data
	*data[tbn] = append(*data[tbn], hb.NewHTMLTableRow(newrow))
	//}
	//create table using the Data and return the HTML
	table := hb.NewHTMLTable(tbn, titles[tbn], *data[tbn], []string{})
	c.JSON(http.StatusOK, responseTableJSON{table.String()})
}
