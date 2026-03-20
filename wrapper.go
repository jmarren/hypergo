package hypergo

import "github.com/a-h/templ"

// type Wrapper func(rw *RW, component templ.Component) templ.Component

type WrapCatcher func(rw *RW, component templ.Component, err error) (templ.Component, error)

type WrapFunc func(rw *RW, component templ.Component) (templ.Component, error)

type Wrapper interface {
	wrap(rw *RW, component templ.Component) templ.Component
	catch(rw *RW, component templ.Component, err error) templ.Component
	Catch(w WrapCatcher) Wrapper
	Wrap(w WrapFunc) Wrapper
}

type wrapper struct {
	wrapFunc   WrapFunc
	catchFuncs []WrapCatcher
}

func (w *wrapper) wrap(rw *RW, component templ.Component) templ.Component {
	component, err := w.wrapFunc(rw, component)

	if err != nil {
		return w.catch(rw, component, err)
	}

	return component

}

func (w *wrapper) catch(rw *RW, component templ.Component, err error) templ.Component {
	for _, catcher := range w.catchFuncs {
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

func (w *wrapper) Catch(catcher WrapCatcher) Wrapper {
	w.catchFuncs = append(w.catchFuncs, catcher)
	return w
}

func (w *wrapper) Wrap(wrapFunc WrapFunc) Wrapper {
	w.wrapFunc = wrapFunc
	return w
}

func defaultWrap(rw *RW, component templ.Component) (templ.Component, error) {
	return component, nil
}

func newWrapper() *wrapper {
	return &wrapper{
		wrapFunc:   defaultWrap,
		catchFuncs: []WrapCatcher{},
	}
}

func UnsafeWrapFunc(fn func(component templ.Component) templ.Component) WrapFunc {
	return func(rw *RW, component templ.Component) (templ.Component, error) {
		return fn(component), nil
	}
}

func UnsafeWrapper(fn func(component templ.Component) templ.Component) Wrapper {
	return SimpleWrapper(UnsafeWrapFunc(fn))
}

func SimpleWrapper(fn WrapFunc) Wrapper {
	w := newWrapper()
	w.Wrap(fn)
	return w
}
