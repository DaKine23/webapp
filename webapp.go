package main

import (
	_ "expvar"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func configMonitoring() {
	http.ListenAndServe(":1234", nil)
}

var switchingValue string

type do func()

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

func main() {
	configRuntime()

	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	//router.LoadHTMLGlob("templates/*.html")
	router.Static("/", "public")

	go everyXSecondsDo(1, switchValue)

	//add cors
	// addCors(router)

	handler := func(c echo.Context) error {
		m := ResponseJSON{switchingValue}
		return c.JSON(http.StatusOK, m)
	}

	router.GET("/pong", handler)
	//GET http://localhost:8080

	nameParamHandler := func(c echo.Context) error {
		name := c.Param("name")

		return c.String(http.StatusOK, "Hello "+name)
	}

	router.GET("/hello/:name", nameParamHandler)
	router.File("favicon.ico", "public/favicon.ico")
	go configMonitoring()

	router.Logger.Fatal(router.Start(":8080"))

}

// ResponseJSON provides a stucture for a message response
type ResponseJSON struct {
	Message string `json:"message" xml:"message"`
}
