package middlewares

import (
	"github.com/gin-generator/ginctl/package/respond"
	"github.com/gin-gonic/gin"
)

// ForceUA Force the request header to include a User-Agent.
func ForceUA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.UserAgent() == "" {
			respond.Alert400WithoutMessage(c, respond.MissUserAgent)
		}
		c.Next()
	}
}
