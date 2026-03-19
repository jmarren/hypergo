package hypergo

import (
	"net/http"

	"github.com/a-h/templ"
)

type Handler func(rw *RW) error
type ComponentHandler func(rw *RW) templ.Component
type Middleware func(h Handler) Handler

type Route struct {
	Parent           *Router
	Path             string
	Method           string
	Handler          Handler
	Middleware       []Middleware
	IsComponent      bool
	ComponentHandler ComponentHandler
	Target           string
}

func makeRWHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &RW{
			ResponseWriter: w,
			Request:        r,
		}
		h(rw)
	}
}

func (route *Route) ComponentHTTPHandler() http.HandlerFunc {
	handler := func(rw *RW) error {

		rw.target = route.Target

		// invoke the componentHandler
		component := route.ComponentHandler(rw)

		// wrap the component up the tree
		parent := route.Parent

		for parent != nil {
			if parent.ShouldWrapPrefix {
				component = parent.PrefixWrapper()(rw, component)
			} else {
				component = parent.Wrapper(rw, component)
			}
			parent = parent.Parent
		}

		if rw.target != "" {
			rw.ResponseWriter.Header().Set("HX-Retarget", rw.target)
		}
		// render & return err
		return component.Render(rw.Request.Context(), rw.ResponseWriter)
	}

	// apply all middleware to the handler
	for _, m := range route.AllMiddleware() {
		handler = m(handler)
	}

	// return a func that creates rw and invokes the handler with it
	return makeRWHandler(handler)
}

func (route *Route) AllMiddleware() []Middleware {
	middleware := route.Middleware

	parent := route.Parent

	for parent != nil {
		middleware = append(middleware, parent.Middleware...)
		parent = parent.Parent
	}
	return middleware
}

func (route *Route) HTTPHandler() http.HandlerFunc {
	handler := route.Handler

	for _, m := range route.AllMiddleware() {
		handler = m(handler)
	}

	return makeRWHandler(handler)

}

func (route *Route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}
