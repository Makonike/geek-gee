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
	r.Run(":9999")
}
