package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// @TODO pass w as a reference!

// Request contains the url parameters and the request body parsed as JSON
type Request struct {
	Parameters      map[string]string
	OriginalRequest *http.Request
	Body            []byte
}

// RouteHandler handles the route
type RouteHandler = func(request Request, w http.ResponseWriter) (interface{}, error)

// Route a route which can be handled
type Route struct {
	Path    string
	Method  string
	Handler RouteHandler
	Regex   string
	Parts   []string
}

// NewRoute creates a new route instance
func NewRoute(path string, method string, handler RouteHandler) *Route {
	route := &Route{Path: path, Method: method, Handler: handler}
	route.Parse(path)
	return route
}

// Match check if the route matches the current url and method
func (r *Route) Match(path string, method string) bool {
	if method != r.Method {
		return false
	}

	res, err := regexp.MatchString(r.Regex, path)
	if err != nil {
		panic(err)
	}
	log.Println("Matching regex", r.Regex, path, res)
	return res
}

// Parameters parse url parameters
func (r *Route) Parameters(path string) map[string]string {
	parts := strings.Split(path, "/")[1:]
	parameters := map[string]string{}
	for key, value := range r.Parts {
		if value == "" {
			continue
		}
		if strings.Contains(value, ":") {
			parameters[strings.Replace(value, ":", "", 1)] = parts[key]
		}
	}
	return parameters
}

// Parse parses a route and generates a regex to match the route
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
	r.Regex = regex + "$"
}

// HTTPServer to describe the httpServer
type HTTPServer struct {
	Routes []*Route
}

func writeResponse(data interface{}, w http.ResponseWriter) {
	if data == nil {
		w.WriteHeader(404)
		return
	}
	p, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(p)
}

func (e *HTTPServer) findRoute(r *http.Request) *Route {
	for _, route := range e.Routes {
		if route.Match(r.URL.Path, r.Method) {
			return route
		}
	}
	return nil
}

func throw500(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)
	w.WriteHeader(500)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Server error"))
}

func throw404(w http.ResponseWriter, r *http.Request) {
	log.Println("error no match for route", r.URL.Path)
	w.WriteHeader(404)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Page not found"))
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	switch err.(type) {
	case *ServerError:
		throw500(err, w, r)
	case *NotFoundError:
		throw404(w, r)
	default:
		throw500(err, w, r)
	}
	return
}

func (e *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting route for url", r.URL.Path, r.Method)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PATCH, POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return
	}

	matchedRoute := e.findRoute(r)
	if matchedRoute == nil {
		throw404(w, r)
		return
	}

	var body []byte
	var err error

	// handle post or patch
	if matchedRoute.Method == http.MethodPost || matchedRoute.Method == http.MethodPatch {
		body, err = ioutil.ReadAll(r.Body)

		if err != nil {
			throw500(err, w, r)
			return
		}
	}

	response, err := matchedRoute.Handler(Request{
		Body:            body,
		Parameters:      matchedRoute.Parameters(r.URL.Path),
		OriginalRequest: r,
	}, w)

	if err != nil {
		handleError(err, w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	writeResponse(response, w)
	return
}

// Register registers a new route
func (e *HTTPServer) Register(path string, method string, handler RouteHandler) {
	e.Routes = append(e.Routes, NewRoute(path, method, handler))
}

// Get register a new route listening to get
func (e *HTTPServer) Get(route string, handler RouteHandler) {
	e.Register(route, http.MethodGet, handler)
}

// Post register a new route listening to post
func (e *HTTPServer) Post(route string, handler RouteHandler) {
	e.Register(route, http.MethodPost, handler)
}

// Patch register a new route listening to patch
func (e *HTTPServer) Patch(route string, handler RouteHandler) {
	e.Register(route, http.MethodPatch, handler)
}

// Delete register a new route listening to delete
func (e *HTTPServer) Delete(route string, handler RouteHandler) {
	e.Register(route, http.MethodDelete, handler)
}

// Listen Listen to a specific port
func (e *HTTPServer) Listen(port string) error {
	return http.ListenAndServe(port, e)
}

type ServerError struct {
	Path string
}

func (e *ServerError) Error() string {
	return "Server error"
}

type UnprocessableEntity struct{}

func (e *UnprocessableEntity) Error() string {
	return "Unprocessable entity"
}

type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return "Resource not found"
}
