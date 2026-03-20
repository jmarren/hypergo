package hypergo

import "github.com/a-h/templ"

func SimpleHandler(fn func() templ.Component) ComponentHandler {
	return func(rw *RW) templ.Component {
		return fn()
	}
}
