package main

import (
	"fmt"
	"net/http"

	"github.com/jmarren/hypergo"
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

func UsernameHandler(rw *hypergo.RW) error {
	rw.ResponseWriter.Write([]byte("username"))
	return nil
}

func RegisterRouter(mux *http.ServeMux, router *hypergo.Router) {

	for _, route := range router.Routes {
		mux.Handle(route.FullPath(), route.HTTPHandler())
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

func main() {

	app := &hypergo.HyperGo{
		Router: &hypergo.Router{
			Path:       "",
			Middleware: []hypergo.Middleware{LoggerOne, LoggerTwo},
			Routes: []*hypergo.Route{
				{
					Path:       "/hi",
					Method:     "GET",
					Handler:    HiHandler,
					Middleware: []hypergo.Middleware{HiHandlerMiddleware},
				},
			},
			SubRouters: []*hypergo.Router{
				{
					Path:       "/users",
					Middleware: []hypergo.Middleware{UsersMiddleware},
					Routes: []*hypergo.Route{
						{
							Path:    "",
							Method:  "GET",
							Handler: UsersHandler,
						},
					},
					SubRouters: []*hypergo.Router{
						{
							Path:       "/username",
							Middleware: []hypergo.Middleware{},
							Routes: []*hypergo.Route{
								{
									Path:    "",
									Method:  "GET",
									Handler: UsernameHandler,
								},
							},
							SubRouters: []*hypergo.Router{},
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
