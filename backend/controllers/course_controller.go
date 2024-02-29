package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Course struct {
	Name string `json:"name" bson:"name"`
}

type CourseController struct {
	Collection *mongo.Collection
}

func NewCourseController(client *mongo.Client) *CourseController {
	return &CourseController{
		Collection: client.Database("ListofCourse").Collection("details"),
	}
}

func (cc *CourseController) UploadCourse(c *gin.Context) {
	var course Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := cc.Collection.InsertOne(context.TODO(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store course data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course data uploaded successfully"})
}

func (cc *CourseController) FetchCourses(c *gin.Context) {
	var courses []Course

	cursor, err := cc.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course data from database"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var course Course
		if err := cursor.Decode(&course); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode course data"})
			return
		}
		courses = append(courses, course)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, courses)
}
