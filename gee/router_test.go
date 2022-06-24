package gee

import (
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/b/c", nil)
	r.addRoute("GET", "/d/:name", nil)
	r.addRoute("GET", "/assert/*filepath", nil)
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/b"), []string{"p", "b"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/makonike")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("assert error!" + n.pattern + " expected: /hello/:name")
	}
	if ps["name"] != "makonike" {
		t.Fatal("assert error!" + ps["name"] + " expected: makonike")
	}

	n, ps = r.getRoute("GET", "/b/c/d")
	if n != nil {
		t.Fatal("nil should be returned")
	}

	n, ps = r.getRoute("GET", "/b")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.pattern != "/b" {
		t.Fatal("assert error!" + n.pattern + " expected: /b")
	}

	n, ps = r.getRoute("GET", "/assert/www")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.pattern != "/assert/*filepath" {
		t.Fatal("assert error!" + n.pattern + " expected: /assert/*filepath")
	}
	if ps["filepath"] != "www" {
		t.Fatal("assert error!" + ps["name"] + " expected: makonike")
	}
}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/b/c", nil)
	r.addRoute("GET", "/b", nil)
	r.addRoute("GET", "/assert/*filepath", nil)
	return r
}
