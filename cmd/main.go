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

type Post struct {
    Title string
    Body  string
}

type Posts []Post

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

func newPost(title string, body string) Post {
    return Post{title, body}
}

func newPosts() Posts {
    return Posts{
		newPost("Title 1", "Body 1"),
		newPost("Title 2", "Body 2"),
		newPost("Title 3", "Body 3"),
	}
}

type Data struct {
    Posts Posts
}

func newData() Data {
    return Data{newPosts()}
}

func main () {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Renderer = newTemplate()
    data := newData()

    e.Static("/css", "css")

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", data)
    })

    e.POST("/markdown", func(c echo.Context) error {
	text := c.FormValue("body")
	parsed := markdown.ToHTML([]byte(text), nil, nil)
	post := newPost("Title", string(parsed))
	return c.HTML(200, post.Body)
	// return c.Render(200, "markdown-content", string(parsed))
    })

    e.POST("/add-post", func(c echo.Context) error {
	body := c.FormValue("body")
	title := c.FormValue("title")
	data.Posts = append(data.Posts, newPost(title, body))
	return c.Render(200, "post-list", data)
    })

    e.Logger.Fatal(e.Start(":3000"))
}

