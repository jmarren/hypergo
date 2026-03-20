package hypergo

import (
	"net/http"

	"github.com/a-h/templ"
)

type Router struct {
	Parent     *Router
	Path       string
	Routes     []*Route
	Middleware []Middleware
	SubRouters []*Router
	Wrapper    Wrapper
	Target     string
}

func defaultWrap(rw *RW, component templ.Component) templ.Component {
	return component
}

func NewRouter(target string) *Router {
	return &Router{
		Parent:     nil,
		Path:       "/",
		Routes:     []*Route{},
		Middleware: []Middleware{},
		SubRouters: []*Router{},
		Target:     target,
		Wrapper:    defaultWrap,
	}
}

func (router *Router) Wrap(w Wrapper) {
	router.Wrapper = func(rw *RW, component templ.Component) templ.Component {
		rw.Retarget(router.Target)
		return w(rw, component)
	}
}

func (router *Router) addRoute(path string, method string, c ComponentHandler) *Route {
	route := &Route{
		Parent:           router,
		Path:             path,
		Method:           method,
		Target:           router.Target,
		Middleware:       router.Middleware,
		ComponentHandler: c,
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

func (router *Router) Get(path string, c ComponentHandler) *Router {
	router.addRoute(path, "GET", c)
	return router
}

func (router *Router) Post(path string, c ComponentHandler) *Router {
	router.addRoute(path, "POST", c)
	return router
}

func (router *Router) Delete(path string, c ComponentHandler) *Router {
	router.addRoute(path, "DELETE", c)
	return router
}

func (router *Router) Put(path string, c ComponentHandler) *Router {
	router.addRoute(path, "PUT", c)
	return router
}

func (router *Router) Patch(path string, c ComponentHandler) *Router {
	router.addRoute(path, "PATCH", c)
	return router
}

func (router *Router) SubRouter(path string, subrouter *Router) {
	subrouter.Parent = router
	subrouter.Path = path
	subrouter.Target = router.Target
	for _, route := range subrouter.Routes {
		route.Middleware = append(route.Middleware, router.Middleware...)
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
		mux.Handle(route.FullPath(), route.ComponentHTTPHandler())
	}

}
