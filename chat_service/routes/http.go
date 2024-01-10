package routers

import (
	"github.com/alireza/handlers"
	"github.com/alireza/middlewares"
	"github.com/gin-gonic/gin"
)

func HTTP(r *gin.Engine) {
	upload := r.Group("/upload", middlewares.Authorize)
	upload.POST("/public", handlers.UploadPublic)
	upload.POST("/private")

}
