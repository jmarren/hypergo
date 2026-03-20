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

func (route *Route) Use(m Middleware) {
	route.Middleware = append([]Middleware{m}, route.Middleware...)
}

// base, users, ...
func (route *Route) ancestors() []*Router {
	ancestors := []*Router{}

	curr := route.Parent

	for curr != nil {
		ancestors = append([]*Router{curr}, ancestors...)
		curr = curr.Parent
	}

	return ancestors
}

func (route *Route) Wrappers(currentPath string) []Wrapper {
	wrappers := []Wrapper{}
	found := false

	ancestors := route.ancestors()

	for _, router := range ancestors {
		currentPath, found = strings.CutPrefix(currentPath, router.Path)
		if !found || currentPath == "" {
			wrappers = append([]Wrapper{router.Wrapper}, wrappers...)
		}
	}

	return wrappers

}

func (route *Route) ComponentHTTPHandler() http.HandlerFunc {
	handler := func(rw *RW) error {
		rw.target = route.Target

		// invoke the componentHandler
		component := route.ComponentHandler(rw)

		wrappers := route.Wrappers(rw.CurrentUrl().Path)

		for _, wrapper := range wrappers {
			component = wrapper(rw, component)
		}

		if rw.target != "" {
			rw.ResponseWriter.Header().Set("HX-Retarget", rw.target)
		}
		// render & return err
		return component.Render(rw.Request.Context(), rw.ResponseWriter)
	}

	// apply all middleware to the handler
	for _, m := range route.Middleware {
		handler = m(handler)
	}

	// return a func that creates rw and invokes the handler with it
	return makeRWHandler(handler, route.Target)
}

func (route *Route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}
