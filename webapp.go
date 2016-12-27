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
	data = []*hb.HTMLTableRow{}
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

var data []*hb.HTMLTableRow
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

	tableResultType := hb.JSONResult{
		Names: []hb.JSONResultName{{"table"}},
	}
	titles := []string{"ping or pong", "timestamp", "id", "delete"}
	tableHandler := func(c *gin.Context) {
		newrow := []interface{}{}
		id++
		newrow = append(newrow, switchingValue)
		newrow = append(newrow, time.Now())
		newrow = append(newrow, id)
		ids := strconv.Itoa(id)

		deletebutton := hb.NewHTMLPart("deletebutton", "", "")

		script := hb.NewScript("tablebutton"+ids, "click", "tablecontainer", "DELETE", "datatable/delete/"+ids, tableResultType, tableResultType.Names[0].Value())
		button := hb.NewHTMLPart("button", "tablebutton"+ids, "del")
		button.AddOption(&hb.HTMLOption{
			Name:  "class",
			Value: "btn btn-danger",
		})
		deletebutton.AddSubPart(script.HTMLPart)
		deletebutton.AddSubPart(button)
		newrow = append(newrow, deletebutton.String())
		data = append(data, hb.NewHTMLTableRow(newrow))

		table := hb.NewHTMLTable("mytable", "tablecontainer", titles, data, []string{})
		c.JSON(http.StatusOK, responseTableJSON{table.String()})
	}
	router.POST("/datatable/add", tableHandler)

	tableDeleteHandler := func(c *gin.Context) {

		id := c.Param("id")

		for i, v := range data {
			if fmt.Sprint((*v.Row)[2]) == id {
				data = append(data[:i], data[i+1:]...)
				break
			}
		}

		table := hb.NewHTMLTable("mytable", "tablecontainer", titles, data, []string{})
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

				data = hb.Sort(i, data)

				break
			}

		}

		table := hb.NewHTMLTable(tbn, "tablecontainer", titles, data, []string{})
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

	html := hb.NewHTMLPart("html", "", "")
	head := hb.NewHTMLPart("head", "", `<meta charset="utf-8">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">`)
	//<link rel="stylesheet" href="static/base.css" />
	//)
	head.AddSubPart(hb.NewHTMLPart("title", "", "Webapp Example"))
	// head.AddSubPart(hb.NewCSSStyle(`table {width: 95%;}
	// 	th {
	// 		background-color: #666; color: #fff;
	// 	}
	// 	tr {
	// 		background-color: #fffbf0; color: #000;
	// 	}
	// 	tr:nth-child(odd) {
	// 		background-color: #e4ebf2 ;
	// 	}
	// 	#tablecontainer tr:hover {
	// 		background-color: #ccc;
	// 	}`))
	jsLibraries := []string{

		"https://unpkg.com/babel-standalone@6.15.0/babel.min.js",
		"https://unpkg.com/jquery@3.1.0/dist/jquery.min.js",
	}

	for _, v := range jsLibraries {
		part := hb.NewHTMLPart("script", "", "")
		part.AddOption(&hb.HTMLOption{
			Name:  "src",
			Value: v,
		})
		head.AddSubPart(part)
	}

	body := hb.NewHTMLPart("body", "", "")
	div := hb.NewHTMLPart("div", "root", "")

	source := "button1"
	destination := "drawdestination"
	source2 := "button2"
	tablecontainer := "tablecontainer"

	resultType := hb.JSONResult{
		Names: []hb.JSONResultName{{"message"}, {"time"}},
	}

	tableResultType := hb.JSONResult{
		Names: []hb.JSONResultName{{"table"}},
	}

	script := hb.NewScript(source, "click", destination, "GET", "/pong", resultType, resultType.Names[0].Value()+`+" !!! " +`+resultType.Names[1].Value())
	scriptToAddARow := hb.NewScript(source2, "click", tablecontainer, "POST", "/datatable/add", tableResultType, tableResultType.Names[0].Value())

	script.AddOption(&hb.HTMLOption{
		Name:  "type",
		Value: "text/babel",
	})

	rows := []*hb.HTMLTableRow{}

	table := hb.NewHTMLTable("mytable", tablecontainer, []string{"ping or pong", "timestamp"}, rows, []string{})

	tp := hb.NewHTMLPart("mytable", tablecontainer, table.String())

	html.AddSubPart(head)
	html.AddSubPart(body)

	body.AddSubPart(div)
	body.AddSubPart(script.HTMLPart)
	body.AddSubPart(scriptToAddARow.HTMLPart)

	button := hb.NewHTMLPart("button", source2, "Add Table Entry")
	button.AddOption(&hb.HTMLOption{
		Name:  "class",
		Value: "btn btn-primary",
	})
	body.AddSubPart(button)

	body.AddSubPart(tp)

	button2 := hb.NewHTMLPart("button", source, "Click Me")
	button2.AddOption(&hb.HTMLOption{
		Name:  "class",
		Value: "btn btn-default",
	})

	body.AddSubPart(button2)
	div2 := hb.NewHTMLPart("div", destination, "content")

	body.AddSubPart(div2)

	return result + html.String()

}
