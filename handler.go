package hypergo

import "github.com/a-h/templ"

func SimpleHandler(fn func() templ.Component) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
		return fn(), nil
	}
}

func SimpleComponent(fn func() templ.Component) Component {
	component := newComponent(SimpleHandler(fn))
	return component

}
