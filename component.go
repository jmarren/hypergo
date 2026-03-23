package hypergo

import (
	"github.com/a-h/templ"
)

type ComponentHandler func(rw *RW) (templ.Component, error)
type ComponentErrCatcher func(rw *RW, component templ.Component, err error) (templ.Component, error)

type Component interface {
	handle(rw *RW) templ.Component
	Catch(catcher ComponentErrCatcher) Component
}

type component struct {
	handler  ComponentHandler
	catchers []ComponentErrCatcher
}

func NewComponent(handler ComponentHandler) *component {
	return &component{
		handler:  handler,
		catchers: []ComponentErrCatcher{},
	}
}

func (c *component) Catch(catcher ComponentErrCatcher) Component {
	c.catchers = append(c.catchers, catcher)
	return c
}

func (c *component) handle(rw *RW) templ.Component {
	component, err := c.handler(rw)

	for _, catcher := range c.catchers {
		component, err = catcher(rw, component, err)
		if err == nil {
			return component
		}
	}

	if err != nil {
		panic(err)
	}
	return component

}
