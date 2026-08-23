package g

import (
	_ "embed"
	"errors"
	"microservice/config"
	"microservice/pkg/router"
	"net/http"
)

//go:embed version
var Version string

//go:embed name
var Name string

var (
	Query     router.ContextKey = "query"
	UserKey   router.ContextKey = "user"
	UuidRegex string            = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
)

// Handling section
type Handler struct {
	Handler func(w http.ResponseWriter, r *http.Request)
}

// Function that gets executed to host a url
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := errors.New("method not implemented")
	if h.Handler != nil {
		h.Handler(w, r)
	} else {
		panic(err)
	}
}

// Config
var CFG *config.Config = nil

// App
var Server *http.Server = nil
