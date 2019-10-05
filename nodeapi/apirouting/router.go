package apirouting

import (
	"net/http"
	"meshed/nodeapi/apilogging"

	"github.com/gorilla/mux"
)

// Returns a preconfigured router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = apilogging.HTTPLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
