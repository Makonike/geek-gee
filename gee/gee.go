package gee

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件，一个分组可以有多个中间件
	parent      *RouterGroup
	engine      *Engine
}

// Engine is the uni handler for all request
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

func New() *Engine {
	engine := &Engine{}
	router := newRouter(engine)
	engine.router = router
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	// add to engine
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute add route
// @param comp: just a part of pattern
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use add middlewares to group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP 解析请求路径，分发处理方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//var middlewares []HandlerFunc
	//// TODO: 每次都要重复遍历太耗费性能了。考虑使用缓存
	//// gin是在前缀树的节点中添加中间件的切片，这样在匹配动态路由并解析参数时，就可以同时获得各分组的中间件
	//for _, v := range engine.groups {
	//	// DONE: FIX /v22, /v23都会匹配到/v2
	//	// 无法实现模糊匹配
	//	if strings.HasPrefix(req.URL.Path, v.prefix) {
	//		middlewares = append(middlewares, v.middlewares...)
	//	}
	//}
	c := newContext(w, req)
	//c.handlers = middlewares
	engine.router.handle(c)
}
