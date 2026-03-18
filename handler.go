package hypergo

//
// type HandlerFunc func(h IHandler)
//
// type Middleware func(h IHandler) IHandler
//
// type IHandler interface {
// 	Req() *http.Request
// 	Res() http.ResponseWriter
// 	isComponent() bool
// }
//
// type Handler struct {
// 	*http.Request
// 	http.ResponseWriter
// 	_isComponent bool
// }
//
// func (h *Handler) Req() *http.Request {
// 	return h.Request
// }
//
// func (h *Handler) Res() http.ResponseWriter {
// 	return h.ResponseWriter
// }
//
// func (h *Handler) isComponent() bool {
// 	return h._isComponent
// }
//
// type ComponentHandler struct {
// 	Handler
// 	component templ.Component
// }
//
// func (c *Handler) IsHxRequest() bool {
// 	return c.Request.Header.Get("HX-Request") == "true"
// }
//
// func (c *ComponentHandler) Render() {
// 	c.component.Render(c.Context(), c.ResponseWriter)
// }
//
// func (c *ComponentHandler) isComponent() bool {
// 	return true
// }
