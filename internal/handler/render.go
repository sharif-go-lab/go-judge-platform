package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// render wraps HTML rendering to inject currentUser
func render(c *gin.Context, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	if cu, ok := c.Get("currentUser"); ok {
		data["currentUser"] = cu
	}
	c.HTML(http.StatusOK, name, data)
}
