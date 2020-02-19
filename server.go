package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	livereload "github.com/mattn/echo-livereload"
)

func main() {
	e := echo.New()
	// logger
	e.Use(middleware.Logger())
	// Stream recovery
	e.Use(middleware.Recover())
	// Live reload
	e.Use(livereload.LiveReload())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// Route => handler
	// Static file handler
	e.Static("/", "assets")
	// html handler
	e.File("/", "ui/index.html")
	// Route => handler
	// Server
	e.Logger.Fatal(e.Start(":8000"))
}
