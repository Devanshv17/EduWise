package routers

import (
	"https://github.com/Devanshv17/StudHelp-IITK/backend/controllers.git"

	"github.com/gin-gonic/gin"
)

func SetupFacultyRouter(r *gin.Engine, facultyController *controllers.FacultyController) {
	api := r.Group("/api")
	{
		api.GET("/faculty", facultyController.FetchFaculty)
		api.POST("/faculty", facultyController.UploadFaculty)
	}
}
