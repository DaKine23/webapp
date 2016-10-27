package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

var switchingValue string

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

func main() {
	configRuntime()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./templates/static")

	go everyXSecondsDo(1, switchValue)

	//add cors
	addCors(router)

	handler := func(c *gin.Context) {
		//c.String(http.StatusOK,"Hello World")

		c.JSON(http.StatusOK, gin.H{
			"message": switchingValue,
		})
	}

	router.GET("/pong", handler)
	//GET http://localhost:3000

	nameParamHandler := func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	}

	router.GET("/hello/:name", nameParamHandler)

	router.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Home Page",
		},
		)
	})
	//GET http://localhost:3000/hello/Rene

	// example requests
	// router.GET("/someGet", getting)
	// router.POST("/somePost", posting)
	// router.PUT("/somePut", putting)
	// router.DELETE("/someDelete", deleting)
	// router.PATCH("/somePatch", patching)
	// router.HEAD("/someHead", head)
	// router.OPTIONS("/someOptions", options)

	router.StaticFile("/favicon.ico", "./favicon.ico")

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
