package main

import (
	"net/http"
	"strconv"

	"github.com/Johnhi19/TreeSpotter_backend/db"
	"github.com/Johnhi19/TreeSpotter_backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	router := gin.Default()

	router.GET("/meadows/:id", findMeadowByID)
	router.POST("/meadows", insertMeadow)

	router.Run("localhost:8080")
}

func insertMeadow(c *gin.Context) {
	var meadow models.Meadow

	// Parse the JSON body into the Meadow struct
	if err := c.ShouldBindJSON(&meadow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the meadow into the database
	insertedID := db.InsertOneMeadow(meadow)

	// Respond with a success message and the new ID
	c.JSON(http.StatusCreated, gin.H{
		"message": "Meadow inserted successfully",
		"id":      insertedID,
	})
}

func findMeadowByID(c *gin.Context) {
	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.D{{"ID", intID}}
	meadow := db.FindOneMeadowById(filter)
	c.IndentedJSON(http.StatusOK, meadow)
}
