package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Faculty struct {
	Name string `json:"name" bson:"name"`
}

type FacultyController struct {
	Collection *mongo.Collection
}

func NewFacultyController(client *mongo.Client) *FacultyController {
	return &FacultyController{
		Collection: client.Database("FacultyInfo").Collection("details"),
	}
}

func (fc *FacultyController) UploadFaculty(c *gin.Context) {
	var faculty Faculty
	if err := c.ShouldBindJSON(&faculty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := fc.Collection.InsertOne(context.TODO(), faculty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store faculty data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty data uploaded successfully"})
}

func (fc *FacultyController) FetchFaculty(c *gin.Context) {
	var faculties []Faculty

	cursor, err := fc.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch faculty data from database"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var faculty Faculty
		if err := cursor.Decode(&faculty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode faculty data"})
			return
		}
		faculties = append(faculties, faculty)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, faculties)
}
