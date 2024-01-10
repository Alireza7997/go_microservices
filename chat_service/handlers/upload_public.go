package handlers

import (
	"github.com/alireza/database"
	"github.com/alireza/pkg/file_manager"
	"github.com/gin-gonic/gin"
)

func UploadPublic(ctx *gin.Context) {
	// database instance
	database := database.Database

	fm := file_manager.New()

	file, _ := ctx.FormFile("file")

}
