package hypergo

import "net/http"

type componentRoute struct {
	*route
	Component ComponentFunc
}

func (route *componentRoute) Handler() http.HandlerFunc {
	handler := func(rw *RW) {
		rw.target = route.Target

		// invoke the componentHandler
		component := route.Component(rw)

		wrappers := route.Wrappers(rw.CurrentUrl().Path)

		for _, wrapper := range wrappers {
			component = wrapper(rw, component)
		}

		if rw.target != "" {
			rw.ResponseWriter.Header().Set("HX-Retarget", rw.target)
		}
		// render & return err
		component.Render(rw.Request.Context(), rw.ResponseWriter)
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
