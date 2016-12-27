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
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	cors "github.com/itsjamie/gin-cors"
	ginglog "github.com/zalando/gin-glog"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

func main() {
	flag.Parse()
	data = []*hb.HtmlTableRow{}
	configRuntime()
	go configMonitoring()

	configController()
}

// ResponseJSON provides a stucture for a message response
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

var data []*hb.HtmlTableRow
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
	router.Static("/static", "./public/static")

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
		c.Writer.WriteString(Page())

		//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
		//c.HTML(200, "index.html", gin.H{})
	})

	tableResultType := hb.JsResult{
		Names: []hb.JsResultName{{"table"}},
	}
	titles := []string{"ping or pong", "timestamp", "id", "delete"}
	tableHandler := func(c *gin.Context) {
		newrow := []interface{}{}
		id++
		newrow = append(newrow, switchingValue)
		newrow = append(newrow, time.Now())
		newrow = append(newrow, id)
		ids := strconv.Itoa(id)

		deletebutton := hb.NewPart("deletebutton", "", "")

		script := hb.NewScript("tablebutton"+ids, "click", "tablecontainer", "DELETE", "datatable/delete/"+ids, tableResultType, tableResultType.Names[0].Value())
		button := hb.NewPart("button", "tablebutton"+ids, "del")
		deletebutton.AddSubPart(script.HtmlPart)
		deletebutton.AddSubPart(button)
		newrow = append(newrow, deletebutton.String())
		data = append(data, hb.NewTableRow(newrow))

		table := hb.NewHtmlTable("mytable", "tablecontainer", titles, data, []string{})
		c.JSON(http.StatusOK, responseTableJSON{table.String()})
	}
	router.POST("/datatable", tableHandler)

	tableDeleteHandler := func(c *gin.Context) {

		id := c.Param("id")

		for i, v := range data {
			if fmt.Sprint((*v.Row)[2]) == id {
				data = append(data[:i], data[i+1:]...)
				break
			}
		}

		table := hb.NewHtmlTable("mytable", "tablecontainer", titles, data, []string{})
		c.JSON(http.StatusOK, responseTableJSON{table.String()})
	}

	router.DELETE("/datatable/delete/:id", tableDeleteHandler)

	newSortHandler := func(c *gin.Context) {
		tbn := c.Param("tablename")
		tbc := c.Param("column")

		for i, v := range titles {
			reducedTitle := strings.Replace(fmt.Sprint(v), " ", "", -1)
			reducedTitle = strings.ToLower(reducedTitle)

			if tbc == reducedTitle {
				sortme := hb.Sorter{}
				sortme.Data = data
				sortme.Sort(i)

				data = sortme.Data

				break
			}

		}

		table := hb.NewHtmlTable(tbn, "tablecontainer", titles, data, []string{})
		c.JSON(http.StatusOK, responseTableJSON{table.String()})
	}

	router.GET("/table/:tablename/:column/sort", newSortHandler)

	//GET http://localhost:3000

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

func Page() string {

	result := "<!DOCTYPE html>\n"

	html := hb.NewPart("html", "", "")
	head := hb.NewPart("head", "", `<meta charset="utf-8">`)
	//<link rel="stylesheet" href="static/base.css" />
	//<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">)
	head.AddSubPart(hb.NewPart("title", "", "Webapp Example"))
	head.AddSubPart(hb.NewCSSStyle(`table {width: 95%;} 
		th { 
			background-color: #666; color: #fff; 
		} 
		tr { 
			background-color: #fffbf0; color: #000; 
		} 
		tr:nth-child(odd) {
			background-color: #e4ebf2 ; 
		}
		#tablecontainer tr:hover { 
   			background-color: #ccc;
		}`))
	jsLibraries := []string{

		"https://unpkg.com/babel-standalone@6.15.0/babel.min.js",
		"https://unpkg.com/jquery@3.1.0/dist/jquery.min.js",
	}

	for _, v := range jsLibraries {
		part := hb.NewPart("script", "", "")
		part.AddOption(&hb.HtmlOption{
			Name:  "src",
			Value: v,
		})
		head.AddSubPart(part)
	}

	body := hb.NewPart("body", "", "")
	div := hb.NewPart("div", "root", "")

	source := "button1"
	destination := "drawdestination"
	source2 := "button2"
	destination2 := "tablecontainer"

	resultType := hb.JsResult{
		Names: []hb.JsResultName{{"message"}, {"time"}},
	}

	tableResultType := hb.JsResult{
		Names: []hb.JsResultName{{"table"}},
	}

	script := hb.NewScript(source, "click", destination, "GET", "/pong", resultType, resultType.Names[0].Value()+`+" !!! " +`+resultType.Names[1].Value())
	script2 := hb.NewScript(source2, "click", destination2, "POST", "/datatable", tableResultType, tableResultType.Names[0].Value())

	///datatable/sort
	script.AddOption(&hb.HtmlOption{
		Name:  "type",
		Value: "text/babel",
	})

	rows := []*hb.HtmlTableRow{}

	table := hb.NewHtmlTable("mytable", "tablecontainer", []string{"ping or pong", "timestamp"}, rows, []string{})

	tp := hb.NewPart("mytable", "tablecontainer", table.String())

	html.AddSubPart(head)
	html.AddSubPart(body)

	body.AddSubPart(div)
	body.AddSubPart(script.HtmlPart)
	body.AddSubPart(script2.HtmlPart)

	body.AddSubPart(hb.NewPart("button", source2, "Add Table Entry"))

	body.AddSubPart(tp)
	body.AddSubPart(hb.NewPart("button", source, "Click Me"))
	div2 := hb.NewPart("div", destination, "content")

	body.AddSubPart(div2)

	return result + html.String()

}
