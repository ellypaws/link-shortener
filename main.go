package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"link-shortener/db"
	"link-shortener/links"
	"net/http"
	"path/filepath"
)

type App struct {
	echo *echo.Echo
	db   *db.Sqlite
}

func NewApp() *App {
	database, err := db.New(context.Background())
	if err != nil {
		panic(err)
	}
	return &App{
		echo: echo.New(),
		db:   database,
	}
}

func main() {
	app := NewApp()

	app.echo.Use(middleware.Logger())
	app.echo.Use(app.storeDatabase)

	// set cors
	app.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	app.echo.GET("/", index)
	app.echo.GET("/:short", retrieveLink)
	app.echo.POST("/shorten", shortenLink)

	app.echo.Logger.Fatal(app.echo.Start(":1323"))
}

func (app App) storeDatabase(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("database", app.db)
		return next(c)
	}
}

func index(c echo.Context) error {
	return c.File("app/dist/index.html")
}

func returnFile(file string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File("app/dist/" + file)
	}
}

func retrieveLink(c echo.Context) error {
	base := filepath.Base(c.Request().URL.Path)
	switch filepath.Ext(base) {
	case ".css", ".js", ".ico":
		return returnFile(c.Request().URL.Path)(c)
	}

	short := c.Param("short")
	database, ok := c.Get("database").(*db.Sqlite)
	if !ok {
		return c.String(500, "database not found")
	}
	link, err := database.SelectLink(short)
	if err != nil {
		return c.String(404, "Link not found")
	}

	return c.Redirect(http.StatusTemporaryRedirect, link.Original)
}

func shortenLink(c echo.Context) error {
	link := links.ShortLink{
		Original: c.FormValue("original"),
		Short:    c.FormValue("short"),
	}

	database, ok := c.Get("database").(*db.Sqlite)
	if !ok {
		return c.String(500, "database not found")
	}

	short := links.NewLink(link.Original, links.WithShort(link.Short))
	err := database.UpsertLink(&short)
	if err != nil {
		return c.String(500, "Failed to add link"+err.Error())
	}

	return c.JSON(200, link)
}
