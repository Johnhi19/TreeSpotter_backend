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
	router.GET("/trees/:id", findTreeByID)
	router.POST("/trees", insertTree)

	router.Run("localhost:8080")
}

func insertMeadow(c *gin.Context) {
	var meadow models.Meadow

	if err := c.ShouldBindJSON(&meadow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedID := db.InsertOneMeadow(meadow)

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

	filter := bson.D{{Key: "ID", Value: intID}}
	meadow := db.FindOneMeadowById(filter)
	c.IndentedJSON(http.StatusOK, meadow)
}

func insertTree(c *gin.Context) {
	var tree models.Tree

	if err := c.ShouldBindJSON(&tree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedID := db.InsertOneTree(tree)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tree inserted successfully",
		"id":      insertedID,
	})
}

func findTreeByID(c *gin.Context) {
	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.D{{Key: "ID", Value: intID}}
	tree := db.FindOneTree(filter)
	c.IndentedJSON(http.StatusOK, tree)
}
