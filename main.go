package main

import (
	"geek-gin/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		//fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
		res := make([]gee.H, 0)
		for k, v := range c.Req.Header {
			res = append(res, gee.H{k: v})
		}
		c.JSON(http.StatusOK, res)
	})
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
