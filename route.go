package hypergo

import "net/http"

type Handler func(rw *RW) error

type Middleware func(h Handler) Handler

type Route struct {
	Parent     *Router
	Path       string
	Method     string
	Handler    Handler
	Middleware []Middleware
}

func (route *Route) HTTPHandler() http.HandlerFunc {
	handler := route.Handler

	for _, m := range route.Middleware {
		handler = m(handler)
	}

	parent := route.Parent

	for parent != nil {
		for _, m := range parent.Middleware {
			handler = m(handler)
		}

		parent = parent.Parent
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(&RW{
			ResponseWriter: w,
			Request:        r,
		})

		if err != nil {
			panic(err)
		}
	}

}

func (route *Route) FullPath() string {
	return route.Method + " " + route.Parent.FullPath() + route.Path
}

// func (route *Route) Register(mux *http.ServeMux) {
// 	mux.Handle(route.method+" "+route.path, route.Handler())
// }
//
// func (route *Route) IHandler(w http.ResponseWriter, r *http.Request) IHandler {
//
// 	if route.isComponent {
// 		return &ComponentHandler{
// 			Handler: Handler{
// 				Request:        r,
// 				ResponseWriter: w,
// 			},
// 		}
// 	}
//
// 	return &Handler{
// 		Request:        r,
// 		ResponseWriter: w,
// 	}
//
// }
