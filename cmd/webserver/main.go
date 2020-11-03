package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	go func() {
		e1 := echo.New()
		e1.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
		e1.Logger.Fatal(e1.Start(":9090"))
	}()

	go func() {
		e2 := echo.New()
		e2.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
		e2.Logger.Fatal(e2.Start(":9091"))
	}()

	select {}
}
