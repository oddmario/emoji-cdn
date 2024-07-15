package main

import (
	"mario/emoji-cdn/routes"

	"github.com/gin-gonic/gin"
)

func initRoutes(engine *gin.Engine) {
	engine.GET("/:emoji", routes.Emoji)
}
