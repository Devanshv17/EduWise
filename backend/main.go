package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://devanshv22:StudHelp@cluster0.adahnkt.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("your_database_name")
	collection = db.Collection("your_collection_name")

	r := gin.Default()

	// Enable CORS
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

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func fetchFiles(c *gin.Context) {
	// Define a struct to represent the data model
	type File struct {
		CourseName string `json:"courseName" bson:"courseName"`
		Year       string `json:"year" bson:"year"`
		Photo      string `json:"photo" bson:"photo"`
	}

	// Define a slice to store fetched data
	var files []File

	// Query MongoDB to fetch data
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from database"})
		return
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode documents into the File struct
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

	// Return fetched data as JSON response
	c.JSON(http.StatusOK, files)
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the file to disk
	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Retrieve other form data
	courseName := c.PostForm("courseName")
	year := c.PostForm("year")

	// Store course info and file path in MongoDB
	_, err = collection.InsertOne(ctx, bson.M{"courseName": courseName, "year": year, "photo": filename})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store data in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful"})
}
