package bootstrap

import (
	"github.com/gin-generator/ginctl/middleware"
	"github.com/gin-gonic/gin"
	demo "github.com/gin-generator/ginctl/app/http/demo/route"
)

func RegisterGlobalMiddleware(r *gin.Engine) {
	r.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.Cors(),
		middlewares.ForceUA(),
	)
}

func RegisterDemoApiRoute(router *gin.Engine) {
	// route not found.
	http.Alert404Route(router)
	// global middleware.
	RegisterGlobalMiddleware(router)
	// Initialize route.
	demo.RegisterDemoAPI(router)
}
