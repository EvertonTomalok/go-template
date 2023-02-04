package rest

import (
	"net/http"

	"github.com/EvertonTomalok/go-template/internal/ui/rest/handlers"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

var healthCheck = []Route{
	{
		"/health",
		http.MethodGet,
		handlers.Health,
	},
	{
		"/readiness",
		http.MethodGet,
		handlers.Readiness,
	},
}

var routes = []Route{
	{
		"/person/:personId",
		http.MethodGet,
		handlers.GetPersonById,
	},
}
