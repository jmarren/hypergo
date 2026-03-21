package hypergo

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type Handler func(rw *RW) error
type ComponentHandler func(rw *RW) (templ.Component, error)

type ComponentErrCatcher func(rw *RW, component templ.Component, err error) (templ.Component, error)
type Component interface {
	handle(rw *RW) templ.Component
	Catch(catcher ComponentErrCatcher) Component
}

type component struct {
	handler  ComponentHandler
	catchers []ComponentErrCatcher
}

func NewComponent(handler ComponentHandler) *component {
	return &component{
		handler:  handler,
		catchers: []ComponentErrCatcher{},
	}
}

func (c *component) Catch(catcher ComponentErrCatcher) Component {
	c.catchers = append(c.catchers, catcher)
	return c
}

func (c *component) handle(rw *RW) templ.Component {
	component, err := c.handler(rw)

	for _, catcher := range c.catchers {
		component, err = catcher(rw, component, err)
		if err == nil {
			return component
		}
	}

	if err != nil {
		panic(err)
	}
	return component

}

type Middleware func(h Handler) Handler

type Route interface {
	Handler() http.HandlerFunc
	Use(m Middleware)
	PrependMiddleware(m ...Middleware)
	FullPath() string
}

type route struct {
	Parent     *Router
	Path       string
	Method     string
	Middleware []Middleware
	Target     string
}

type componentRoute struct {
	*route
	Component Component
}

type regularRoute struct {
	*route
	handler Handler
	catcher func(rw *RW, err error)
}

func (route *regularRoute) Handler() http.HandlerFunc {

	handler := route.handler
	// apply all middleware to the handler
	for _, m := range route.Middleware {
		handler = m(handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rw := &RW{
			ResponseWriter: w,
			Request:        r,
		}
		err := handler(rw)

		if err != nil {
			route.catcher(rw, err)
		}
	}

}

func (route *route) Use(m Middleware) {
	route.Middleware = append([]Middleware{m}, route.Middleware...)
}

func (route *route) PrependMiddleware(m ...Middleware) {
	route.Middleware = append(route.Middleware, m...)
}

// base, users, ...
func (route *route) ancestors() []*Router {
	ancestors := []*Router{}

	curr := route.Parent

	for curr != nil {
		ancestors = append([]*Router{curr}, ancestors...)
		curr = curr.Parent
	}

	return ancestors
}

func (route *route) Wrappers(currentPath string) []Wrapper {
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

func (route *componentRoute) Handler() http.HandlerFunc {
	handler := func(rw *RW) error {
		rw.target = route.Target

		// invoke the componentHandler
		component := route.Component.handle(rw)

		wrappers := route.Wrappers(rw.CurrentUrl().Path)

		for _, wrapper := range wrappers {
			component = wrapper.wrap(rw, component)
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

	return func(w http.ResponseWriter, r *http.Request) {
		rw := newRW(w, r, route.Target)
		handler(rw)
	}

}

func (route *route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}
