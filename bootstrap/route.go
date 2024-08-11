package bootstrap

import (
	admin "github.com/gin-generator/ginctl/app/http/admin/route"
	"github.com/gin-generator/ginctl/middleware"
	http "github.com/gin-generator/ginctl/package/respond"
	"github.com/gin-gonic/gin"
)

func RegisterGlobalMiddleware(r *gin.Engine) {
	r.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.Cors(),
		middlewares.ForceUA(),
	)
}

func RegisterAdminApiRoute(router *gin.Engine) {
	// route not found.
	http.Alert404Route(router)
	// global middleware.
	RegisterGlobalMiddleware(router)
	// Initialize route.
	admin.RegisterAdminAPI(router)
}
