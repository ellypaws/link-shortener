package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"link-shortener/links"
	"net/http"
)

type App struct {
	Links *links.Linker
	echo  *echo.Echo
}

func NewApp() *App {
	return &App{
		Links: links.NewLinker(),
		echo:  echo.New(),
	}
}

func main() {
	app := NewApp()

	app.echo.Use(middleware.Logger())
	app.echo.Use(app.storeLinker)

	// set cors
	app.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	app.echo.GET("/:short", retrieveLink)
	app.echo.POST("/shorten", shortenLink)

	app.echo.Logger.Fatal(app.echo.Start(":1323"))
}

func (app App) storeLinker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("linker", app.Links)
		return next(c)
	}
}

func retrieveLink(c echo.Context) error {
	short := c.Param("short")
	linker, ok := c.Get("linker").(*links.Linker)
	if !ok {
		return c.String(500, "linker not found")
	}
	link := linker.GetLink(short)
	if link == nil {
		return c.String(404, "Link not found")
	}

	return c.Redirect(http.StatusTemporaryRedirect, link.Original)
}

func shortenLink(c echo.Context) error {
	link := links.ShortLink{
		Original: c.FormValue("original"),
		Short:    c.FormValue("short"),
	}

	linker, ok := c.Get("linker").(*links.Linker)
	if !ok {
		return c.String(500, "linker not found")
	}

	short := linker.NewLink(link.Original, links.WithShort(link.Short))
	err := linker.AddLink(short)
	if err != nil {
		return c.String(500, "Failed to add link")
	}

	return c.JSON(200, link)
}
