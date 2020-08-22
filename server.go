package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/imthaghost/musik/soundcloud"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

var songlist = [...]string{"https://soundcloud.com/uiceheidd/tell-me-you-love-me", "https://soundcloud.com/gunna/top-floor-feat-travis-scott", "https://soundcloud.com/shallou/forget", "https://soundcloud.com/liltjay/zoo-york-feat-fivio-foreign",
	"https://soundcloud.com/polo-g/polo-g-feat-juice-wrld-flex", "https://soundcloud.com/polo-g/polo-g-feat-lil-baby-be", "https://soundcloud.com/uiceheidd/righteous", "https://soundcloud.com/meduzamusic/lose-control", "https://soundcloud.com/tenth999/juice-wrld-boys-in-da-band",
	"https://soundcloud.com/roddyricch/the-box", "https://soundcloud.com/lil-baby-4pf/sum-2-prove", "https://soundcloud.com/ianndior/prospect-feat-lil-baby", "https://soundcloud.com/poorchoice/chateau", "https://soundcloud.com/uiceheidd/blood-on-my-jeans"}
var old string

func main() {
	path := "assets/music"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	e := echo.New()

	// Log Output
	e.Use(middleware.Logger())
	// stream
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

		rand.Seed(time.Now().Unix())
		var n = rand.Int() % len(songlist)

		g := songlist[n]
		songname, image, path := soundcloud.ExtractSong(g)
		old = path
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"url":     g,
			"name":    songname,
			"artwork": image,
			"song":    "music/" + path,
		})
	}).Name = "index"

	e.Logger.Fatal(e.Start(":5000"))
}
