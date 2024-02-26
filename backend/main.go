package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var facultyCollection *mongo.Collection
var courseCollection *mongo.Collection
var ctx = context.TODO()

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

type Faculty struct {
	Name string `json:"name" bson:"name"`
}

type Course struct {
	Name string `json:"name" bson:"name"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://devanshv22:StudHelp@cluster0.adahnkt.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("CoursesInfo")
	collection = db.Collection("details")

	facultyDB := client.Database("FacultyInfo")
	facultyCollection = facultyDB.Collection("details")

	courseDB := client.Database("ListofCourse")
	courseCollection = courseDB.Collection("details")

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.POST("/api/upload", uploadFile)
	r.GET("/api/fetch", fetchFiles)
	r.GET("/api/faculty", fetchFaculty)
	r.POST("/api/faculty", uploadFaculty)
	r.POST("/api/courses", uploadCourse)
	r.GET("/api/courses", fetchCourses)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func generateFileLink(fileID string) string {
	// Assuming your frontend is hosted at http://localhost:3000
	return fmt.Sprintf("http://localhost:3000/files/%s", fileID)
}

func uploadFile(c *gin.Context) {
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
	_, err = collection.InsertOne(ctx, bson.M{
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

func fetchFiles(c *gin.Context) {
	var files []File

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from database"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var file File
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

func uploadFaculty(c *gin.Context) {
	var faculty Faculty
	if err := c.ShouldBindJSON(&faculty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := facultyCollection.InsertOne(ctx, faculty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store faculty data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty data uploaded successfully"})
}

func fetchFaculty(c *gin.Context) {
	var faculties []Faculty

	cursor, err := facultyCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch faculty data from database"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
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

func uploadCourse(c *gin.Context) {
	var course Course // Assuming you have a struct definition for Course similar to Faculty
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := courseCollection.InsertOne(ctx, course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store course data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course data uploaded successfully"})
}

func fetchCourses(c *gin.Context) {
	var courses []Course // Assuming you have a struct definition for Course similar to Faculty

	cursor, err := courseCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course data from database"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
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
