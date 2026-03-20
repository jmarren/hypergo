package main

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jmarren/hypergo"
	"github.com/jmarren/hypergo/views"
)

func LoggerOne(h hypergo.Handler) hypergo.Handler {
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

func WrapPage(rw *hypergo.RW, component templ.Component) templ.Component {
	username, ok := rw.Context().Value("username").(string)
	if !ok {
		fmt.Printf("username not found in ctx\n")
	}
	return views.Base(views.Page(component, username))
}

func main() {

	app := hypergo.New("#content")

	app.Use(LogRequest)
	app.Use(LoggerOne)
	app.Use(LoggerTwo)
	app.Use(AddUsername)
	app.Wrap(WrapPage)

	usersRouter := hypergo.NewRouter("#users-component")
	usersRouter.Wrap(hypergo.SimpleWrapper(views.Users))

	usersRouter.Get("username", hypergo.SimpleHandler(views.Username))
	usersRouter.Use(LoggerThree)
	usersRouter.Get("age", hypergo.SimpleHandler(views.Age))

	songsRouter := hypergo.NewRouter("#songs-component")

	songsRouter.Wrap(hypergo.SimpleWrapper(views.Songs))

	songsRouter.Get("blackbird", hypergo.SimpleHandler(views.Blackbird))

	YesterdayRouter := hypergo.NewRouter("#yesterday-component")

	YesterdayRouter.Wrap(hypergo.SimpleWrapper(views.Yesterday))

	YesterdayRouter.Get("stats", hypergo.SimpleHandler(views.YesterdayStats))
	YesterdayRouter.Get("artwork", hypergo.SimpleHandler(views.YesterdayArtwork))

	songsRouter.SubRouter("yesterday/", YesterdayRouter)

	app.Router.SubRouter("users/", usersRouter)
	app.Router.SubRouter("songs/", songsRouter)

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
