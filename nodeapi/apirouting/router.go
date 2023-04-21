package apirouting

import (
	"github.com/gorilla/mux"
	"github.com/mgrote/meshed/nodeapi/apilogging"
	"net/http"
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
