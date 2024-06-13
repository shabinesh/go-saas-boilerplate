package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handlers) HomePage(r *gin.Context) {
	r.HTML(http.StatusOK, "home", gin.H{})
}
