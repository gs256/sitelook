package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeRoute(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/search")
}
