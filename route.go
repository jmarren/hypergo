package hypergo

import (
	"net/http"
	"strings"
)

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

func (route *route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}
