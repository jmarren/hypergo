package hypergo

import (
	"net/http"

	"github.com/a-h/templ"
)

type Router struct {
	Parent     *Router
	Path       string
	Routes     []Route
	Middleware []Middleware
	SubRouters []*Router
	Wrapper    Wrapper
	Target     string
}

func NewRouter(target string) *Router {
	return &Router{
		Parent:     nil,
		Path:       "/",
		Routes:     []Route{},
		Middleware: []Middleware{},
		SubRouters: []*Router{},
		Target:     target,
		Wrapper:    newWrapper(),
	}
}

func (router *Router) Wrap(w WrapFunc) Wrapper {
	router.Wrapper.Wrap(func(rw *RW, component templ.Component) (templ.Component, error) {
		rw.Retarget(router.Target)
		return w(rw, component)
	})
	return router.Wrapper

}

func (router *Router) addComponentRoute(path string, method string, c Component) Route {
	route := &componentRoute{
		route: &route{
			Parent:     router,
			Path:       path,
			Method:     method,
			Target:     router.Target,
			Middleware: router.Middleware,
		},
		Component: c,
	}
	router.Routes = append(router.Routes, route)
	return route

}

func (router *Router) addRegularRoute(path string, method string, h Handler) Route {
	route := &regularRoute{
		route: &route{
			Parent:     router,
			Path:       path,
			Method:     method,
			Target:     router.Target,
			Middleware: router.Middleware,
		},
		handler: h,
	}
	router.Routes = append(router.Routes, route)
	return route

}

func (router *Router) Use(m Middleware) {
	router.Middleware = append([]Middleware{m}, router.Middleware...)
}

func (router *Router) SetTarget(target string) {
	router.Target = target
}

func (router *Router) GetComponent(path string, c Component) Component {
	router.addComponentRoute(path, "GET", c)
	return c
}

func (router *Router) Get(path string, h Handler) Handler {
	router.addRegularRoute(path, "GET", h)
	return h
}

func (router *Router) Post(path string, h Handler) Handler {
	router.addRegularRoute(path, "POST", h)
	return h
}

func (router *Router) Delete(path string, c Component) Component {
	router.addComponentRoute(path, "DELETE", c)
	return c
}

func (router *Router) Put(path string, c Component) Component {
	router.addComponentRoute(path, "PUT", c)
	return c
}

func (router *Router) Patch(path string, c Component) Component {
	router.addComponentRoute(path, "PATCH", c)
	return c
}

func (router *Router) SubRouter(path string, subrouter *Router) {
	subrouter.Parent = router
	subrouter.Path = path
	subrouter.Target = router.Target
	for _, route := range subrouter.Routes {
		route.PrependMiddleware(router.Middleware...)
		// route.Middleware = append(route.Middleware, router.Middleware...)
		router.Routes = append(router.Routes, route)
	}
}

func (router *Router) FullPath() string {
	path := ""
	parent := router

	for parent != nil {
		path = parent.Path + path
		parent = parent.Parent
	}

	return path
}

func (router *Router) register(mux *http.ServeMux) {

	for _, route := range router.Routes {
		mux.Handle(route.FullPath(), route.Handler())
	}

}
