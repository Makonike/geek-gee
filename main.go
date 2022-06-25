package main

import (
	"fmt"
	"geek-gin/gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		t := time.Now()
		// 使用v2符合route的请求都会走这里,否则会走默认的404handler
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/", func(c *gee.Context) {
		//fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
		res := make([]gee.H, 0)
		for k, v := range c.Req.Header {
			res = append(res, gee.H{k: v})
		}
		c.JSON(http.StatusOK, res)
	})
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("D:\\1\\nginx\\nginx-1.20.1\\course_selection\\html\\html\\account\\user/*")
	r.Static("/css", "D:\\1\\nginx\\nginx-1.20.1\\course_selection\\html\\css/")
	r.Static("/images", "D:\\1\\nginx\\nginx-1.20.1\\course_selection\\html\\images/")
	r.GET("/ab", func(c *gee.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
	}
	r.GET("/hello", func(c *gee.Context) {
		//fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
		res := make([]gee.H, 0)
		for k, v := range c.Req.Header {
			res = append(res, gee.H{k: v})
		}
		c.JSON(http.StatusOK, res)
	})
	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	r.GET("/asserts/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})
	r.Run(":9999")
}
