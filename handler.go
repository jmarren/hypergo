package hypergo

import "github.com/a-h/templ"

// type Handler func(rw *RW) error

//go:generate ./build/gatekeeper -type=Pill
type HandleFunc func(rw *RW) error
type Catcher func(rw *RW, err error) error

type Handler interface {
	handle(rw *RW)
	Catch(catcher ...Catcher) Handler
}

type handler struct {
	handlerFunc HandleFunc
	catchers    []Catcher
	validator   RequestValidator
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

func NewHandler(h HandleFunc) Handler {
	return &handler{
		handlerFunc: h,
		catchers:    []Catcher{},
		validator:   NewRequestValidator(),
	}
}

func SimpleHandler(fn func() templ.Component) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
		return fn(), nil
	}
}

func SimpleComponent(fn func() templ.Component) Component {
	component := NewComponent(SimpleHandler(fn))
	return component
}
