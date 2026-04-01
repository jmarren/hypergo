package main

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jmarren/hypergo"
	"github.com/jmarren/hypergo/pages"
	"github.com/jmarren/hypergo/views"
)

// type Middleware func(h HandlerFunc) HandlerFunc
func LoggerOne(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Printf("loggerOne\n")
		h(rw)
	}
}

func LoggerTwo(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Printf("loggerTwo\n")
		h(rw)
	}
}

func LoggerThree(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Printf("loggerThree\n")
		h(rw)
	}
}

func AddUsername(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Println("adding username")
		rw.Request = rw.Request.WithContext(context.WithValue(rw.Context(), "username", "john"))
		h(rw)
	}
}

func LogRequest(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Printf("%s %s\n", rw.Request.Method, rw.URL.Path)
		h(rw)
	}
}

func WrapBase(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Base(component)
}

func WrapPage(rw *hypergo.RW, component templ.Component) templ.Component {
	username, ok := rw.Context().Value("username").(string)
	if !ok {
		return views.Base(views.Page(component, "user not found"))
	}
	return views.Base(views.Page(component, username))
}

func pageCatcher(rw *hypergo.RW, component templ.Component, err error) (templ.Component, error) {
	if err.Error() == "username not found" {
		return views.Base(views.Page(component, "not found")), nil
	}

	return component, err
}

func main() {

	app := hypergo.New("#content")

	app.Use(LogRequest)
	app.Use(LoggerOne)
	app.Use(LoggerTwo)
	app.HxWrap(WrapPage)

	app.GetComponent("", hypergo.SimpleComponent(views.Age))
	app.GetComponent("about", hypergo.SimpleComponent(views.Blackbird))

	app.Router.SubRouter("users/", pages.UsersRouter)
	app.Router.SubRouter("songs/", pages.SongsRouter)

	app.Listen(":5050")

}
