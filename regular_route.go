package hypergo

import "net/http"

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
