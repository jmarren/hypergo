package hypergo

import (
	"github.com/a-h/templ"
)

type Wrapper func(rw *RW, component templ.Component) templ.Component

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

func NewRouter() *Router {
	return &Router{
		Parent:     nil,
		Path:       "/",
		Routes:     []*Route{},
		Middleware: []Middleware{},
		SubRouters: []*Router{},
		Target:     "",
		Wrapper:    defaultWrap,
	}
}

func (router *Router) Wrap(w Wrapper) {
	router.Wrapper = func(rw *RW, component templ.Component) templ.Component {
		rw.Retarget(router.Target)
		return w(rw, component)
	}
}

func (router *Router) addRoute(path string, method string, target string, c ComponentHandler) *Route {
	route := &Route{
		Parent:           router,
		Path:             path,
		Method:           method,
		Target:           target,
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

func (router *Router) Get(path string, target string, c ComponentHandler) *Router {
	router.addRoute(path, "GET", target, c)
	return router
}

func (router *Router) Post(path string, target string, c ComponentHandler) *Router {
	router.addRoute(path, "POST", target, c)
	return router
}

func (router *Router) Delete(path string, target string, c ComponentHandler) *Router {
	router.addRoute(path, "DELETE", target, c)
	return router
}

func (router *Router) Put(path string, target string, c ComponentHandler) *Router {
	router.addRoute(path, "PUT", target, c)
	return router
}

func (router *Router) Patch(path string, target string, c ComponentHandler) *Router {
	router.addRoute(path, "PATCH", target, c)
	return router
}

func (router *Router) SubRouter(path string, subrouter *Router) {
	subrouter.Parent = router
	subrouter.Path = path
	for _, route := range subrouter.Routes {
		route.Middleware = append(router.Middleware, route.Middleware...)
	}
	router.SubRouters = append(router.SubRouters, subrouter)
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
