package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/RuaanV/styles/1.0.0/",
		Index,
	},

	Route{
		"Index",
		"GET",
		"/RuaanV/skus/1.0.0/",
		Index,
	},

	Route{
		"CreateSkuCommand",
		"POST",
		"/RuaanV/skus/1.0.0/sku/command/createsku",
		CreateSkuCommand,
	},

	Route{
		"CreateStyleCommand",
		"POST",
		"/RuaanV/styles/1.0.0/syle/command/createstyle",
		CreateStyleCommand,
	},

	Route{
		"PingGet",
		"GET",
		"/RuaanV/skus/1.0.0/ping",
		PingGet,
	},

	Route{
		"PingGet",
		"GET",
		"/RuaanV/styles/1.0.0/ping",
		PingGet,
	},

	Route{
		"GetSkuById",
		"GET",
		"/RuaanV/skus/1.0.0/sku/{skuId}",
		GetSkuByIdWithDbConnection(),
	},

	Route{
		"GetStyleById",
		"GET",
		"/RuaanV/styles/1.0.0/style/{styleId}",
		GetStyleByIdWithDbConnection(),
	},

	Route{
		"ListSkus",
		"GET",
		"/RuaanV/skus/1.0.0/skus",
		ListSkus,
	},

	Route{
		"ListStyles",
		"GET",
		"/RuaanV/styles/1.0.0/styles",
		ListStyles,
	},
}
