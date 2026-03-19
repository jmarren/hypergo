package hypergo

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type Handler func(rw *RW) error
type ComponentHandler func(rw *RW) templ.Component
type Middleware func(h Handler) Handler

type Route struct {
	Parent           *Router
	Path             string
	Method           string
	Middleware       []Middleware
	ComponentHandler ComponentHandler
	Target           string
}

func makeRWHandler(h Handler, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := newRW(w, r, target)
		h(rw)
	}
}

func (route *Route) ComponentHTTPHandler() http.HandlerFunc {
	handler := func(rw *RW) error {
		rw.target = route.Target

		// invoke the componentHandler
		component := route.ComponentHandler(rw)

		rw.URL.Path, _ = strings.CutSuffix(rw.URL.Path, route.Path)
		// wrap the component up the tree
		parent := route.Parent

		currentUrl := rw.CurrentUrl().Path

		for rw.URL.Path != currentUrl {
			component = parent.Wrapper(rw, component)
			rw.URL.Path, _ = strings.CutSuffix(rw.URL.Path, parent.Path)
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
	return makeRWHandler(handler, route.Target)
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

func (route *Route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}
