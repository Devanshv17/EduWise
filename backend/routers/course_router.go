package routers

import (
	"github.com/devanshv17/StudHelp-IITK/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCourseRouter(r *gin.Engine, courseController *controllers.CourseController) {
	api := r.Group("/api")
	{
		api.POST("/courses", courseController.UploadCourse)
		api.GET("/courses", courseController.FetchCourses)
	}
}
