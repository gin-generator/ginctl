package bootstrap

import (
	"github.com/gin-generator/ginctl/middleware"
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
