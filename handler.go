package hypergo

import "github.com/a-h/templ"

type HandlerFunc func(rw *RW)
type Catcher func(rw *RW, err error) error

type Handler interface {
	handle(rw *RW)
	Catch(catcher ...Catcher) Handler
	HandleFunc() HandlerFunc
}

type handler struct {
	handlerFunc func(rw *RW) error
	catchers    []Catcher
}

func (h *handler) HandleFunc() HandlerFunc {
	return h.handle
}

func (h *handler) handle(rw *RW) {
	err := h.handlerFunc(rw)

	if err == nil {
		return
	}

	for _, catcher := range h.catchers {
		err = catcher(rw, err)
		if err == nil {
			return
		}
	}
	if err != nil {
		panic(err)
	}
}

func (h *handler) Catch(catchers ...Catcher) Handler {
	h.catchers = append(h.catchers, catchers...)
	return h
}

func NewHandler(h func(rw *RW) error) Handler {
	return &handler{
		handlerFunc: h,
		catchers:    []Catcher{},
	}
}

func SimpleHandler(fn func() templ.Component) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
		return fn(), nil
	}
}

func SimpleComponent(fn func() templ.Component) ComponentFunc {
	return NewComponent(SimpleHandler(fn)).Handler()
}
