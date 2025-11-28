package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Johnhi19/TreeSpotter_backend/db"
	"github.com/Johnhi19/TreeSpotter_backend/models"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect("credentialsMySql.txt")
	defer db.Disconnect()

	router := gin.Default()

	router.GET("/meadows/:id", findMeadowByID)
	router.GET("/meadows", getBasicInfoOfAllMeadows)
	router.POST("/meadows", insertMeadow)
	router.GET("/meadows/:id/trees", getTreesOfMeadow)
	router.GET("/trees/:id", findTreeByID)
	router.POST("/trees", insertTree)

	router.Run("localhost:8080")

	go func() {
		if err := router.Run("localhost:8080"); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	db.Disconnect()
}

func getBasicInfoOfAllMeadows(c *gin.Context) {
	meadows := db.FindAllMeadows()
	c.IndentedJSON(http.StatusOK, meadows)
}

func getTreesOfMeadow(c *gin.Context) {
	meadowId := c.Param("id")

	intID, err := strconv.Atoi(meadowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	trees := db.FindAllTreesForMeadow(intID)
	c.IndentedJSON(http.StatusOK, trees)
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

	meadow := db.FindOneMeadowById(intID)
	c.IndentedJSON(http.StatusOK, meadow)
}

func insertTree(c *gin.Context) {
	var tree models.Tree

	if err := c.ShouldBindJSON(&tree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the tree
	insertedID := db.InsertOneTree(tree)

	// Update the meadow's TreeIds list
	if err := db.UpdateMeadowTreeIds(tree.MeadowId, insertedID); err != nil {
		fmt.Printf("ERROR executing UPDATE: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tree inserted but failed to update meadow"})
		return
	}

	fmt.Printf("Updated Meadow %d with new Tree ID %d\n", tree.MeadowId, insertedID)

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

	tree := db.FindOneTreeById(intID)
	c.IndentedJSON(http.StatusOK, tree)
}
