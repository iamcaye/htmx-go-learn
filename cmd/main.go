package main

import (
    "html/template"
    "io"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    "github.com/gomarkdown/markdown"
)

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}


func main () {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Renderer = newTemplate()

    e.Static("/css", "css")

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", nil)
    })

    e.POST("/markdown", func(c echo.Context) error {
	text := c.FormValue("markdown")
	parsed := markdown.ToHTML([]byte(text), nil, nil)
	return c.HTML(200, string(parsed))
	// return c.Render(200, "markdown-content", string(parsed))
    })

    e.Logger.Fatal(e.Start(":3000"))
}

