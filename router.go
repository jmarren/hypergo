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

func (router *Router) Get(c ComponentHandler) *Route {

	return route
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

// func (router *Router) PrefixWrapper() Wrapper {
// 	return func(rw *RW, component templ.Component) templ.Component {
// 		fullPath := router.FullPath()
// 		fmt.Printf("fullPath = %s\n", fullPath)
//
// 		if rw.PathHasPrefix(fullPath) {
// 			return component
// 		}
// 		rw.Retarget(router.Target)
// 		return router.Wrapper(rw, component)
// 	}
//
// }
//
// func (router *Router) WrapPrefix() {
// 	currWrapper := router.Wrapper
//
// 	router.Wrapper = func(rw *RW, component templ.Component) templ.Component {
// 		fullPath := router.FullPath()
// 		fmt.Printf("fullPath = %s\n", fullPath)
//
// 		if rw.PathHasPrefix(fullPath) {
// 			return component
// 		}
// 		rw.Retarget(router.Target)
// 		return currWrapper(rw, component)
// 	}
// }

// func (r *Router)

// func (r *Router) Use(m Middleware) *Router {
// 	r.middleware = append([]Middleware{m}, r.middleware...)
// 	return r
// }
