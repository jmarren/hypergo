package main

import (
	"fmt"
	"net/http"

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

func HiHandlerMiddleware(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Printf("HiHandlerMiddleware\n")
		return h(rw)
	}
}

func UsersMiddleware(h hypergo.Handler) hypergo.Handler {
	return func(rw *hypergo.RW) error {
		fmt.Printf("users\n")
		return h(rw)
	}
}

func HiHandler(rw *hypergo.RW) error {
	rw.ResponseWriter.Write([]byte("hi"))
	return nil
}

func UsersHandler(rw *hypergo.RW) error {
	rw.ResponseWriter.Write([]byte("users"))
	return nil
}

// func UsernameHandler(rw *hypergo.RW) error {
// 	rw.ResponseWriter.Write([]byte("username"))
// 	return nil
// }

func RegisterRouter(mux *http.ServeMux, router *hypergo.Router) {

	for _, route := range router.Routes {
		if route.IsComponent {
			mux.Handle(route.FullPath(), route.ComponentHTTPHandler())
		} else {
			mux.Handle(route.FullPath(), route.HTTPHandler())
		}
	}

	for _, subrouter := range router.SubRouters {
		RegisterRouter(mux, subrouter)
	}
}

func assignParents(router *hypergo.Router) {
	for _, route := range router.Routes {
		route.Parent = router
	}

	for _, subrouter := range router.SubRouters {
		subrouter.Parent = router
		assignParents(subrouter)
	}
}

func WrapBase(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Base(component)
}

func WrapPage(rw *hypergo.RW, component templ.Component) templ.Component {

	if rw.IsHxRequest() {
		return component
	}
	return views.Base(views.Page(component, "john"))
}

func WrapUsers(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Users(component)
}

func WrapSongs(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Songs(component)
}

func WrapYesterday(rw *hypergo.RW, component templ.Component) templ.Component {
	return views.Yesterday(component)
}

func UsernameHandler(rw *hypergo.RW) templ.Component {
	return views.Username()
}

func YesterdayStatsHandler(rw *hypergo.RW) templ.Component {
	return views.YesterdayStats()
}

func YesterdayArtworkHandler(rw *hypergo.RW) templ.Component {
	return views.YesterdayArtwork()
}

func BlackbirdHandler(rw *hypergo.RW) templ.Component {
	return views.Blackbird()
}

func AgeHandler(rw *hypergo.RW) templ.Component {
	return views.Age()
}

func EmptyWrap(rw *hypergo.RW, component templ.Component) templ.Component {
	return component
}

func main() {

	app := &hypergo.HyperGo{
		Router: &hypergo.Router{
			Wrapper:          WrapPage,
			ShouldWrapPrefix: false,
			Target:           "body",
			Path:             "",
			Middleware:       []hypergo.Middleware{LoggerOne, LoggerTwo},
			Routes: []*hypergo.Route{
				{
					Path:        "/hi",
					Method:      "GET",
					Handler:     HiHandler,
					IsComponent: false,
					Middleware:  []hypergo.Middleware{HiHandlerMiddleware},
				},
			},
			SubRouters: []*hypergo.Router{
				{
					Wrapper:          WrapUsers,
					ShouldWrapPrefix: true,
					Target:           "#content",
					Path:             "/users",
					Middleware:       []hypergo.Middleware{UsersMiddleware},
					Routes: []*hypergo.Route{
						{
							Path:        "",
							Method:      "GET",
							Handler:     UsersHandler,
							IsComponent: false,
						},
						{
							Path:             "/username",
							Method:           "GET",
							ComponentHandler: UsernameHandler,
							Target:           "#users-component",
							IsComponent:      true,
						},
						{
							Path:             "/age",
							Method:           "GET",
							Target:           "#users-component",
							ComponentHandler: AgeHandler,
							IsComponent:      true,
						},
					},
				},
				{
					Wrapper:          WrapSongs,
					ShouldWrapPrefix: true,
					Target:           "#content",
					Path:             "/songs",
					Middleware:       []hypergo.Middleware{UsersMiddleware},
					Routes: []*hypergo.Route{
						{
							Path:             "/blackbird",
							Method:           "GET",
							Target:           "#songs-component",
							ComponentHandler: BlackbirdHandler,
							IsComponent:      true,
						},
					},
					SubRouters: []*hypergo.Router{
						{
							Wrapper:          WrapYesterday,
							Path:             "/yesterday",
							ShouldWrapPrefix: true,
							Target:           "#songs-component",
							Routes: []*hypergo.Route{
								{
									Path:             "/stats",
									Method:           "GET",
									ComponentHandler: YesterdayStatsHandler,
									IsComponent:      true,
									Target:           "#yesterday-component",
								},
								{
									Path:             "/artwork",
									Method:           "GET",
									ComponentHandler: YesterdayArtworkHandler,
									IsComponent:      true,
									Target:           "#yesterday-component",
								},
							},
						},
					},
				},
			},
		},
	}

	mux := http.NewServeMux()

	assignParents(app.Router)

	RegisterRouter(mux, app.Router)

	// mux := http.NewServeMux()

	fmt.Printf("app = %v\n", app)

	http.ListenAndServe(":5050", mux)

}

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
