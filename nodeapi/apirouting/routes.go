package apirouting

import (
	"meshed/nodeapi/apihandler"
	"net/http"
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
		"GET",
		"/node/{" + apihandler.TypeName + "}/{" + apihandler.NodeID + ":[0-9a-z]+}",
		apihandler.ShowNode,
	},
	Route{
		"UploadFile",
		"POST",
		"/upload",
		apihandler.UploadFileHandler,
	},
}
