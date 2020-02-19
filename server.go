package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	livereload "github.com/mattn/echo-livereload"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	// Live Reload
	e.Use(livereload.LiveReload())
	// Log Output
	e.Use(middleware.Logger())
	// Stream Recovery
	e.Use(middleware.Recover())
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// static files
	e.Static("/", "assets")
	// template render
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	// Named route "index"
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name":       "Ransom - Lil Tecca & Juice Wrld",
			"audioURL":   "http://localhost:5000/music/ransom.mp3",
			"artworkURL": "https://content-images.p-cdn.com/images/public/int/7/8/3/3/00602508243387_1080W_1080H.jpg",
		})
	}).Name = "index"

	e.Logger.Fatal(e.Start(":5000"))
}
