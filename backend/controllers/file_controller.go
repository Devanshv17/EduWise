package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type File struct {
	ID          string `json:"id" bson:"_id"`
	CourseName  string `json:"courseName" bson:"courseName"`
	Batch       string `json:"batch" bson:"batch"`
	Instructor  string `json:"instructor" bson:"instructor"`
	Type        string `json:"type" bson:"type"`
	Remark      string `json:"remark" bson:"remark"`
	FileContent []byte `json:"fileContent" bson:"fileContent"`
	Link        string `json:"link" bson:"link"`
}

type FileController struct {
	Collection *mongo.Collection
}

func NewFileController(client *mongo.Client) *FileController {
	return &FileController{
		Collection: client.Database("CoursesInfo").Collection("details"),
	}
}

func (fc *FileController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open uploaded file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer fileContent.Close()

	// Read file content
	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Generate a unique identifier for the file
	fileID := primitive.NewObjectID().Hex()

	courseName := c.PostForm("courseName")
	batch := c.PostForm("batch")
	instructor := c.PostForm("instructor")
	fileType := c.PostForm("type")
	remark := c.PostForm("remark")

	// Store file content and other details in the database
	_, err = fc.Collection.InsertOne(context.TODO(), bson.M{
		"_id":         fileID, // Store the unique identifier
		"courseName":  courseName,
		"batch":       batch,
		"instructor":  instructor,
		"type":        fileType,
		"remark":      remark,
		"fileContent": content,                  // Store file content
		"link":        generateFileLink(fileID), // Store the generated link
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful"})
}

func (fc *FileController) FetchFiles(c *gin.Context) {
	type FileWithoutContent struct {
		ID         string `json:"id" bson:"_id"`
		CourseName string `json:"courseName" bson:"courseName"`
		Batch      string `json:"batch" bson:"batch"`
		Instructor string `json:"instructor" bson:"instructor"`
		Type       string `json:"type" bson:"type"`
		Remark     string `json:"remark" bson:"remark"`
		Link       string `json:"link" bson:"link"`
	}

	var files []FileWithoutContent

	cursor, err := fc.Collection.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"fileContent": 0}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from database"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var file FileWithoutContent
		if err := cursor.Decode(&file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode data"})
			return
		}

		files = append(files, file)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (fc *FileController) DownloadFile(c *gin.Context) {
	fileID := c.Param("fileID")

	// Fetch file details from the database using the fileID
	var file File
	err := fc.Collection.FindOne(context.TODO(), bson.M{"_id": fileID}).Decode(&file)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file content
	c.Data(http.StatusOK, "application/octet-stream", file.FileContent)
}

func generateFileLink(fileID string) string {
	// Assuming your frontend is hosted at http://localhost:3000
	return fmt.Sprintf("http://localhost:8080/api/download/%s?download=true", fileID)
}
