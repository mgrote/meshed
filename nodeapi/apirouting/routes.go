package apirouting

import (
	"github.com/mgrote/meshed/nodeapi/apihandler"
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
		"/nodes/{" + apihandler.NodeTypeName + "}",
		apihandler.ListNodes,
	},
	Route{
		"ShowNode",
		"GET",
		"/node/{" + apihandler.NodeTypeName + "}/{" + apihandler.NodeID + ":[0-9a-z]+}",
		apihandler.ShowNode,
	},
	Route{
		"UploadFile",
		"POST",
		"/upload",
		apihandler.UploadFileHandler,
	},
	Route{
		"RegisterUser",
		"POST",
		"/register",
		apihandler.RegisterUser,
	},
	Route{
		"LoginUser",
		"POST",
		"/login",
		apihandler.LoginUser,
	},
	Route{
		"RenewUserToken",
		"GET",
		"/renew",
		apihandler.RenewUserToken,
	},
}
