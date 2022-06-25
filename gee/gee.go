package gee

import (
	"html/template"
	"net/http"
	"path"
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
	router        *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // 模板标准库
	funcMap       template.FuncMap
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
	c.engine = engine
	//c.handlers = middlewares
	engine.router.handle(c)
}

// SetFuncMap 将所有的模板加载进内存
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 自定义模板渲染函数
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 创建静态文件处理器
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// 从请求 URL 的路径中删除给定的前缀并调用处理程序 http.FileServer(fs)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// Static 解析请求地址并映射到服务器上文件的真实地址
func (group *RouterGroup) Static(relativePath, root string) {
	// TODO: http.Dir()可能会暴露敏感文件和目录
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
