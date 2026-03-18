package hypergo

type Router struct {
	Parent     *Router
	Path       string
	Routes     []*Route
	Middleware []Middleware
	SubRouters []*Router
}

func (r *Router) FullPath() string {
	path := ""
	parent := r

	for parent != nil {
		path = parent.Path + path
		parent = parent.Parent
	}

	return path
}

// func (r *Router) Use(m Middleware) *Router {
// 	r.middleware = append([]Middleware{m}, r.middleware...)
// 	return r
// }
