package hypergo

import "net/http"

type HyperGo struct {
	mux *http.ServeMux
	*Router
}

func New(target string) *HyperGo {
	return &HyperGo{
		Router: NewRouter(target),
		mux:    http.NewServeMux(),
	}
}

func (h *HyperGo) Listen(addr string) {
	h.register(h.mux)
	http.ListenAndServe(addr, h.mux)
}
