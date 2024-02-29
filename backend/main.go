package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	setupRoutes(r)
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	r.Use(setupCors())

	// Initialize MongoDB client
	ctx := context.TODO()
	client, err := NewMongoClient(ctx, "mongodb+srv://devanshv22:StudHelp@cluster0.adahnkt.mongodb.net/?retryWrites=true&w=majority")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Set up routes
	fileController := NewFileController(client)
	facultyController := NewFacultyController(client)
	courseController := NewCourseController(client)

	api := r.Group("/api")
	{
		api.POST("/upload", fileController.UploadFile)
		api.GET("/fetch", fileController.FetchFiles)
		api.GET("/faculty", facultyController.FetchFaculty)
		api.POST("/faculty", facultyController.UploadFaculty)
		api.POST("/courses", courseController.UploadCourse)
		api.GET("/courses", courseController.FetchCourses)
		api.GET("/download/:fileID", fileController.DownloadFile)
	}
}

func setupCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
