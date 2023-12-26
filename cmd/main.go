package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/iamcaye/htmx-go-learn/cmd/models"
	"github.com/iamcaye/htmx-go-learn/cmd/repos"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    templs := &Templates{}

    templs.templates = template.Must(template.New("").Funcs(
	template.FuncMap{
	    "markdown": func(html string) template.HTML {
		    return template.HTML(markdown.ToHTML([]byte(html), nil, nil))
		},
	    "safe": func(html string) template.HTML {
		return template.HTML(html)
	    },
	}).ParseGlob("views/*.html"))

    templs.templates.ParseGlob("views/components/*.html")
    fmt.Println(templs.templates.DefinedTemplates())
    return templs;
}

func newPost(title string, body string) models.Post {
    return models.Post{
	Title: title,
	Body:  body,
    }
}

func newPosts() models.Posts {
    return models.Posts{
		newPost("Title 1", "Body 1"),
		newPost("Title 2", "Body 2"),
		newPost("Title 3", "Body 3"),
	}
}

type Data struct {
    Posts models.Posts
    SelectedPost models.Post
    EditPost bool
}

func newData() Data {
    return Data{
	Posts: newPosts(),
	SelectedPost: models.Post{},
	EditPost: true,
    }
}

func (d *Data) AddPost(post models.Post) {
	d.Posts = append(d.Posts, post)
}

func (d *Data) AddPosts(posts models.Posts) {
    d.Posts = append(d.Posts, posts...)
}

func (d *Data) ParseMarkdown () Data {
    for i := 0; i < len(d.Posts); i++ {
	parsed := markdown.ToHTML([]byte(d.Posts[i].Body), nil, nil)
	d.Posts[i].Body = string(parsed)
    }
    return *d
}

func main () {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Renderer = newTemplate()
    data := newData()
    e.Static("/css", "css")
    e.GET("/", func(c echo.Context) error {
	mongoClient, err := repos.InitMongo()
	if err != nil {
	    fmt.Println(err)
	    return err
	}

	posts, err := repos.GetPosts(mongoClient)
	if err == nil {
	    data.Posts = posts
	}

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
	mongoClient, err := repos.InitMongo()
	if err != nil {
	    fmt.Println(err)
	    return err
	}

	post := newPost(title, body)

	err = repos.AddPost(mongoClient, post)
	if err != nil {
	    fmt.Println(err)
	    return err
	}

	data.Posts = append(data.Posts, post)
	c.Render(200, "add-post-button", nil)
	return c.Render(200, "post-list", data)
    })

    e.GET("/edit-post", func(c echo.Context) error {
	return c.Render(200, "post-form", newPost("", ""))
    })

    e.GET("/cancel-post", func(c echo.Context) error {
	return c.Render(200, "add-post-button", nil)
    })

    e.Logger.Fatal(e.Start(":3000"))
}

