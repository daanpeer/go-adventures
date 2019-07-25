package main

import (
	"fmt"
	"net/http"
)

type Route struct {
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type HttpServer struct {
	Routes map[string]Route
}

func (e *HttpServer) listen() {
	http.ListenAndServe(":8080", e)
}

// could add routing with parameters e.g. /:id/
// automatic response
func (e *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := e.Routes[r.URL.Path]

	fmt.Println(route, r.URL.Path, e.Routes, r.Method, route.Method)

	if ok && r.Method == route.Method {
		route.Handler(w, r)
	}
}

func (e *HttpServer) Register(path string, route Route) {
	if e.Routes == nil {
		e.Routes = make(map[string]Route)
	}
	e.Routes[path] = route
}

func (e *HttpServer) Get(route string, handler func(w http.ResponseWriter, r *http.Request)) {
	e.Register(route, Route{Method: http.MethodGet, Handler: handler})
}

func (e *HttpServer) Post(route string, handler func(w http.ResponseWriter, r *http.Request)) {
	e.Register(route, Route{Method: http.MethodPost, Handler: handler})
}

func main() {
	app := HttpServer{}

	app.Get("/kaas", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("kaas")
	})

	app.Post("/koek", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("koek")
	})

	app.listen()

}
