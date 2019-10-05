package apirouting

import (
	"net/http"
	"meshed/nodeapi/apihandler"
)

// Route def
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Array of Route
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		apihandler.IndexRootHandler,
	},
	Route{
		"ListEntryPoints",
		"GET",
		"/listtypes",
		apihandler.ListNodeTypes,
	},
	Route{
		"ListNodes",
		"GET",
		"/nodes/{" + apihandler.TypeName + "}",
		apihandler.ListNodes,
	},
	Route{
		"ShowNode",
		"POST",
		"/node/{" + apihandler.TypeName + "}/{" + apihandler.NodeID + "}",
		apihandler.ShowNode,
	},
}
