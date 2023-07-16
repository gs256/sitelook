package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "home-page", nil)
}
