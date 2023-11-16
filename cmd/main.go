package main

import (
	"html/template"
	"io"
	"log"

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
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

type Contact struct {
    Name string
    Email string
}

func newContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}

type Contacts []Contact

type Data struct {
    Contacts Contacts
}

func newData() Data {
    return Data{
        Contacts: []Contact{
            newContact("caye", "caye@caye.es"),
            newContact("arnold", "arnold@arnold.es"),
        },
    }
}

func main () {
    e := echo.New()
    e.Use(middleware.Logger())

    data := newData()

    e.Renderer = newTemplate()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", data)
    })

    e.POST("/add-contact", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        contact := newContact(name, email)
        data.Contacts = append(data.Contacts, contact)

        err := c.Render(200, "contact-form", data);
        if (err != nil) {
            log.Fatal(err)
        }
        return c.Render(200, "oob-contact-item", contact)
    })

    e.Logger.Fatal(e.Start(":3000"))
}

