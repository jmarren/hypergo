package hypergo

import "net/http"

type componentRoute struct {
	*route
	Component ComponentFunc
}

func (route *componentRoute) Handler() http.HandlerFunc {
	handler := func(rw *RW) {
		rw.target = route.Target

		// invoke the componentHandler to get the routes component
		component := route.Component(rw)

		// get the neccessary wrappers based on the current url
		wrappers := route.Wrappers(rw.CurrentUrl().Path)

		// wrap the component with all neccessary wrappers
		for _, wrapper := range wrappers {
			component = wrapper(rw, component)
		}

		// if the target is set, retarget the response
		if rw.target != "" {
			rw.ResponseWriter.Header().Set("HX-Retarget", rw.target)
		}
		// render
		component.Render(rw.Request.Context(), rw.ResponseWriter)
	}

	// apply all middleware to the handler
	for _, m := range route.Middleware {
		handler = m(handler)
	}

	// return an http.HandlerFunc that converts r and w to RW and pass it to the handler
	return func(w http.ResponseWriter, r *http.Request) {
		rw := newRW(w, r, route.Target)
		handler(rw)
	}

}
