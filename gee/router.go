package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路径
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			// 模糊匹配剪枝，匹配的路径到此，剩下会作为参数
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 添加路由
func (r *router) addRoute(method, pattern string, handlerFunc HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handlerFunc
	log.Printf("Route: %4s - %s", method, pattern)
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	// not this route
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, p := range parts {
			if p[0] == ':' {
				params[p[1:]] = searchParts[i]
			}
			if p[0] == '*' && len(p) > 1 {
				params[p[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 将匹配到的handler都加紧到执行handler中，调用Next()去执行
func (r *router) handle(c *Context) {
	n, param := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = param
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		// 默认的404handler
		c.handlers = append(c.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}
	c.Next()
}
