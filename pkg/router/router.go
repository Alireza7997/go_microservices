package router

import (
	"context"
	"fmt"
	"microservice/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Context type and keys
type ContextKey string

func ReturnContextKey(input string) ContextKey {
	return ContextKey(input)
}

type muxEntry struct {
	h         map[string]http.Handler
	pattern   string
	variables map[int]variable
	methods   []string
}

type variableType string

type variable struct {
	t    variableType
	name string
}

const (
	variableString variableType = "string"
	variableInt    variableType = "int64"
	variableFloat  variableType = "float64"
)

type Router struct {
	http.ServeMux
	mu          sync.RWMutex
	m           map[string]muxEntry
	es          []muxEntry
	middlewares []func(http.Handler) http.Handler
}

var allowed_methods = []string{
	"GET", "POST", "PUT", "PATCH", "DELETE",
}

func convertToType(in string) variableType {
	if in == "string" {
		return variableString
	} else if in == "int" {
		return variableInt
	} else if in == "float" {
		return variableFloat
	}

	return variableString
}

// Adds a middleware to the whole process of running
// Every request
func (mux *Router) Middleware(handler func(http.Handler) http.Handler) {
	mux.middlewares = append(mux.middlewares, handler)
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *Router) Handle(pattern string, handler http.Handler, method string) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	method = strings.ToUpper(method)
	allowed := false
	for _, allowed_method := range allowed_methods {
		allowed = allowed_method == method || allowed
	}
	if !allowed {
		panic(fmt.Sprintf("`%s` is not a valid http method", method))
	}

	if pattern == "" || pattern == "//" {
		panic("router: invalid pattern")
	}
	if handler == nil {
		panic("router: nil handler")
	}
	if element, exist := mux.m[pattern]; exist {
		methodExist := false
		for _, m := range element.methods {
			methodExist = methodExist || m == method
		}
		if methodExist {
			panic("router: multiple registrations for " + pattern)
		}
	}
	if pattern[len(pattern)-1] != '/' || pattern[0] != '/' {
		panic("router: first and end of every route has to contain `/`")
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	// Record variable names
	splits := strings.Split(pattern, "/")
	reformatPattern := ""
	variableRegex, _ := regexp.Compile(`^{([a-zA-Z_]\w*)}|{([a-zA-Z_]\w*):(int|float|string)}`)
	variables := map[int]variable{}
	for i, split := range splits {
		if len(split) == 0 {
			continue
		}

		if reformatPattern == "" {
			reformatPattern = "/"
		}

		actualVariable := variableRegex.FindString(split)
		if actualVariable != "" {
			actualVariable = actualVariable[1 : len(actualVariable)-1]
			if splits = strings.Split(actualVariable, ":"); len(splits) > 1 {
				variables[i] = variable{name: splits[0], t: convertToType(splits[1])}
			} else {
				variables[i] = variable{name: actualVariable, t: variableString}
			}
		}

		split = variableRegex.ReplaceAllString(split, "")
		reformatPattern += split + "/"
	}
	if reformatPattern == "" {
		reformatPattern = "/"
	}
	reformatPattern = fmt.Sprintf("^%s$", reformatPattern)

	if element, exist := mux.m[pattern]; exist {
		element.methods = append(element.methods, method)
		element.h[method] = handler
		return
	}
	e := muxEntry{h: map[string]http.Handler{method: handler}, pattern: reformatPattern, variables: variables, methods: []string{method}}
	mux.m[pattern] = e
	mux.es = append(mux.es, e)
}

// Runs the handler with middlewares
func (mux *Router) runWithMiddlewares(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, _, variables := mux.handler(path, r.Method)
		r = mux.recordVariables(path, variables, r)
		h.ServeHTTP(w, r)
	})

	var function http.Handler
	for i := len(mux.middlewares) - 1; i > -1; i-- {
		mid := mux.middlewares[i]
		if function == nil {
			function = mid(handler)
			continue
		}
		function = mid(function)
	}
	function.ServeHTTP(w, r)
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	if path != "/" {
		if path[len(path)-1] != '/' {
			path = fmt.Sprintf("%s/", path)
		}
	}
	r.URL.Path = path

	mux.runWithMiddlewares(w, r)
}

// Records variables inside url into context
func (mux *Router) recordVariables(path string, variables map[int]variable, r *http.Request) *http.Request {
	ctx := r.Context()

	splits := strings.Split(path, "/")
	for i, variable := range variables {
		value := splits[i]
		if variable.t == variableInt {
			actualVariable, _ := strconv.ParseInt(value, 10, 64)
			ctx = context.WithValue(ctx, ReturnContextKey(variable.name), actualVariable)
		} else if variable.t == variableFloat {
			actualVariable, _ := strconv.ParseFloat(value, 64)
			ctx = context.WithValue(ctx, ReturnContextKey(variable.name), actualVariable)
		} else {
			ctx = context.WithValue(ctx, ReturnContextKey(variable.name), value)
		}
	}
	return r.WithContext(ctx)
}

// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
func (mux *Router) handler(path, method string) (h http.Handler, pattern string, variables map[int]variable) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if h == nil {
		atLeastOneRouteMatched := false
		found := false
		for i := 0; i < len(mux.es); i++ {
			element := mux.es[i]
			r, _ := regexp.Compile(element.pattern)
			routeMatched := r.MatchString(path)
			atLeastOneRouteMatched = routeMatched || atLeastOneRouteMatched
			_, methodMatched := element.h[method]
			if routeMatched && methodMatched {
				h = element.h[method]
				pattern = element.pattern
				variables = element.variables
				found = true
				break
			}
		}
		if atLeastOneRouteMatched && !found {
			panic(errors.New(http.StatusMethodNotAllowed, "", "method not allowed"))
		} else if !found {
			panic(errors.New(http.StatusNotFound, "page not found", ""))
		}
	}
	return
}
