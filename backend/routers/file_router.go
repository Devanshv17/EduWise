package routers

import (
	"https://github.com/Devanshv17/StudHelp-IITK/backend/controllers.git"

	"github.com/gin-gonic/gin"
)

func SetupFileRouter(r *gin.Engine, fileController *controllers.FileController) {
	api := r.Group("/api")
	{
		api.POST("/upload", fileController.UploadFile)
		api.GET("/fetch", fileController.FetchFiles)
		api.GET("/download/:fileID", fileController.DownloadFile)
	}
}
