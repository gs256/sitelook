package app

import (
	"sitelook/app/home"
	"sitelook/app/search"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	engine := gin.Default()

	engine.GET("/", home.HomeRoute)
	// engine.GET("/api/search", apiSearchRoute)
	engine.GET("/search", search.SearchRoute)

	engine.Static("./static", "./static/")
	engine.LoadHTMLGlob("templates/*")
	engine.Run()
}
