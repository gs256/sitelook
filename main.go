package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func searchRoute(c *gin.Context) {
	searchTerm := c.Query("q")
	startQuery := c.Query("start")

	if len(searchTerm) == 0 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString("empty search term")
		return
	}

	start, _ := strconv.Atoi(startQuery)
	searchPageContext, err := parseSearchPage(searchTerm, start)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, searchPageContext)
}

func main() {
	r := gin.Default()
	r.GET("/api/search", searchRoute)
	r.Run()
}
