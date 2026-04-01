package hypergo

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentFunc func(rw *RW) templ.Component

type ComponentWrapper func(rw *RW, c templ.Component) templ.Component

type Router struct {
	Parent     *Router
	Path       string
	Routes     []Route
	Middleware []Middleware
	SubRouters []*Router
	Wrapper    ComponentWrapper
	Target     string
}

func EmptyWrap(rw *RW, c templ.Component) templ.Component {
	return c
}

func NewRouter(target string) *Router {
	return &Router{
		Parent:     nil,
		Path:       "/",
		Routes:     []Route{},
		Middleware: []Middleware{},
		SubRouters: []*Router{},
		Target:     target,
		Wrapper:    EmptyWrap,
	}
}

func (router *Router) HxWrap(w ComponentWrapper) *Router {
	router.Wrapper = func(rw *RW, c templ.Component) templ.Component {
		if rw.IsHxRequest() {
			return c
		}
		rw.Retarget(router.Target)
		return w(rw, c)
	}
	return router
}

func (router *Router) Wrap(w ComponentWrapper) *Router {
	router.Wrapper = func(rw *RW, c templ.Component) templ.Component {
		rw.Retarget(router.Target)
		return w(rw, c)
	}
	return router
}

func (router *Router) addComponentRoute(path string, method string, c ComponentFunc) Route {
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

func (router *Router) addRegularRoute(path string, method string, h HandlerFunc) Route {
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

func (router *Router) GetComponent(path string, c ComponentFunc) ComponentFunc {
	router.addComponentRoute(path, "GET", c)
	return c
}

func (router *Router) Get(path string, h HandlerFunc) *Router {
	router.addRegularRoute(path, "GET", h)
	return router
}

func (router *Router) Post(path string, h HandlerFunc) *Router {
	router.addRegularRoute(path, "POST", h)
	return router
}

func (router *Router) Delete(path string, h HandlerFunc) *Router {
	router.addRegularRoute(path, "DELETE", h)
	return router
}

func (router *Router) Put(path string, h HandlerFunc) *Router {
	router.addRegularRoute(path, "PUT", h)
	return router
}

func (router *Router) Patch(path string, h HandlerFunc) *Router {
	router.addRegularRoute(path, "PATCH", h)
	return router
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
