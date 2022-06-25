package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Recovery 使用defer挂载上错误恢复的函数,返回500
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				ctx.JSON(http.StatusInternalServerError, H{"code": http.StatusInternalServerError})
			}
		}()
		ctx.Next()
	}
}

// 用于触发panic的堆栈信息
func trace(msg string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(msg + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
