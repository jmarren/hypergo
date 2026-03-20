package hypergo

import "github.com/a-h/templ"

type Wrapper func(rw *RW, component templ.Component) templ.Component

func SimpleWrapper(fn func(c templ.Component) templ.Component) Wrapper {
	return func(rw *RW, component templ.Component) templ.Component {
		return fn(component)
	}
}
