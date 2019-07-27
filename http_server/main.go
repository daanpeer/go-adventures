package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	Method  string
	Handler RouteHandler
	Regex   string
	Parts   []string
}

func NewRoute(path string, method string, handler RouteHandler) *Route {
	route := &Route{Method: method, Handler: handler}
	route.Parse(path)
	return route
}

func (r *Route) Match(path string, method string) bool {
	if method != r.Method {
		return false
	}
	res, err := regexp.MatchString(r.Regex, path)
	if err != nil {
		panic(err)
	}
	return res
}

func (r *Route) Parameters(path string) map[string]string {
	parts := strings.Split(path, "/")
	parameters := map[string]string{}
	for key, value := range r.Parts {
		if strings.Contains(value, ":") {
			parameters[strings.Replace(value, ":", "", 1)] = parts[key]
		}
	}
	return parameters
}

func (r *Route) Parse(path string) {
	parts := strings.Split(path, "/")
	r.Parts = parts
	var regex string
	for _, value := range parts {
		if value == "" {
			continue
		}

		if strings.Contains(value, ":") {
			regex += `\/([a-z0-9]*)`
			continue
		}
		regex += fmt.Sprintf(`\/(%s)`, value)
	}
	r.Regex = regex
}

type HttpServer struct {
	Routes map[string]*Route
}

func (e *HttpServer) Listen() {
	http.ListenAndServe(":8080", e)
}

// could add routing with parameters e.g. /:id/
// automatic response
func (e *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var matchedRoute *Route
	for _, route := range e.Routes {
		if route.Match(r.URL.Path, r.Method) {
			matchedRoute = route
			break
		}
	}

	if matchedRoute != nil {
		parameters := matchedRoute.Parameters(r.URL.Path)
		matchedRoute.Handler(parameters, w, r)
		return
	}

	fmt.Println("error no match for route", r.URL.Path)
}

type RouteHandler = func(parameters map[string]string, w http.ResponseWriter, r *http.Request)

func (e *HttpServer) Register(path string, method string, handler RouteHandler) {
	if e.Routes == nil {
		e.Routes = make(map[string]*Route)
	}
	e.Routes[path] = NewRoute(path, method, handler)
}

func (e *HttpServer) Get(route string, handler RouteHandler) {
	e.Register(route, http.MethodGet, handler)
}

func (e *HttpServer) Post(route string, handler RouteHandler) {
	e.Register(route, http.MethodGet, handler)
}

func main() {
	app := HttpServer{}

	app.Get("/kaas/:id", func(parameters map[string]string, w http.ResponseWriter, r *http.Request) {
		fmt.Println("kaas/id", parameters)
	})

	app.Get("/kaas", func(parameters map[string]string, w http.ResponseWriter, r *http.Request) {
		fmt.Println("kaas")
	})

	app.Post("/koek", func(parameters map[string]string, w http.ResponseWriter, r *http.Request) {
		fmt.Println("koek")
	})

	app.Listen()

}
