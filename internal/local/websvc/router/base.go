package router

import (
	"community/internal/local/websvc/handler"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Routers() *gin.Engine {
	Router.Use(gzip.Gzip(gzip.DefaultCompression))

	Router.GET("/health", handler.Handler)
	Router.GET("/add", handler.AddT)

	return Router
}
