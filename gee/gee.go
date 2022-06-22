package gee

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

// Engine is the uni handler for all request
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute add route
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}
func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.router.addRoute("PUT", pattern, handler)
}

func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.router.addRoute("DELETE", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 解析请求路径，分发处理方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
