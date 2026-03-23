package main

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jmarren/hypergo"
	"github.com/jmarren/hypergo/pages"
	"github.com/jmarren/hypergo/views"
)

func LoggerOne(h hypergo.HandleFunc) hypergo.HandleFunc {
	return func(rw *hypergo.RW) error {
		fmt.Printf("loggerOne\n")
		return h(rw)
	}
}

func LoggerTwo(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Printf("loggerTwo\n")
		return h(rw)
	}
}

func LoggerThree(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Printf("loggerThree\n")
		return h(rw)
	}
}

func AddUsername(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Println("adding username")
		rw.Request = rw.Request.WithContext(context.WithValue(rw.Context(), "username", "john"))
		return h(rw)
	}
}

func LogRequest(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Printf("%s %s\n", rw.Request.Method, rw.URL.Path)
		return h(rw)
	}
}

func WrapBase(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Base(component)
}

func WrapPage(rw *hypergo.RW, component templ.Component) (templ.Component, error) {
	username, ok := rw.Context().Value("username").(string)
	if !ok {
		return component, fmt.Errorf("username not found")
	}
	return views.Base(views.Page(component, username)), nil
}

func pageCatcher(rw *hypergo.RW, component templ.Component, err error) (templ.Component, error) {
	if err.Error() == "username not found" {
		return views.Base(views.Page(component, "not found")), nil
	}

	return component, err
}

func main() {

	// hypergo.TryValidate()

	app := hypergo.New("#content")

	app.Use(LogRequest)
	app.Use(LoggerOne)
	app.Use(LoggerTwo)
	app.Wrap(WrapPage).Catch(pageCatcher)

	app.Router.SubRouter("users/", pages.UsersRouter)
	app.Router.SubRouter("songs/", pages.SongsRouter)

	app.Listen(":5050")

}

// app := &hypergo.HyperGo{
// 	Router: &hypergo.Router{
// 		Wrapper:    WrapPage,
// 		Target:     "body",
// 		Path:       "",
// 		Middleware: []hypergo.Middleware{LoggerOne, LoggerTwo},
// 		Routes:     []*hypergo.Route{},
// 		SubRouters: []*hypergo.Router{
// 			{
// 				Wrapper:    WrapUsers,
// 				Target:     "#content",
// 				Path:       "/users",
// 				Middleware: []hypergo.Middleware{UsersMiddleware},
// 				Routes: []*hypergo.Route{
// 					{
// 						Path:             "/username",
// 						Method:           "GET",
// 						ComponentHandler: UsernameHandler,
// 						Target:           "#users-component",
// 					},
// 					{
// 						Path:             "/age",
// 						Method:           "GET",
// 						Target:           "#users-component",
// 						ComponentHandler: AgeHandler,
// 					},
// 				},
// 			},
// 			{
// 				Wrapper:    WrapSongs,
// 				Target:     "#content",
// 				Path:       "/songs",
// 				Middleware: []hypergo.Middleware{UsersMiddleware},
// 				Routes: []*hypergo.Route{
// 					{
// 						Path:             "/blackbird",
// 						Method:           "GET",
// 						Target:           "#songs-component",
// 						ComponentHandler: BlackbirdHandler,
// 					},
// 				},
// 				SubRouters: []*hypergo.Router{
// 					{
// 						Wrapper: WrapYesterday,
// 						Path:    "/yesterday",
// 						Target:  "#songs-component",
// 						Routes: []*hypergo.Route{
// 							{
// 								Path:             "/stats",
// 								Method:           "GET",
// 								ComponentHandler: YesterdayStatsHandler,
// 								Target:           "#yesterday-component",
// 							},
// 							{
// 								Path:             "/artwork",
// 								Method:           "GET",
// 								ComponentHandler: YesterdayArtworkHandler,
// 								Target:           "#yesterday-component",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

// SubRouters: []*hypergo.Router{
// 	{
// 		Path:             "/username",
// 		ShouldWrapPrefix: false,
// 		Middleware:       []hypergo.Middleware{},
// 		Wrapper:          EmptyWrap,
// 		Routes: []*hypergo.Route{
// 			{
// 				Path:             "",
// 				Method:           "GET",
// 				IsComponent:      true,
// 				ComponentHandler: UsernameHandler,
// 			},
// 		},
// 		SubRouters: []*hypergo.Router{},
// 	},
// },
