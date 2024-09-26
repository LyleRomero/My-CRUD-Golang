package main

import (
	"My-CRUD-Golang/internal/adapters/db"
	"My-CRUD-Golang/internal/application"
	"My-CRUD-Golang/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := db.NewMemoryRepository()
	service := application.NewItemService(repo)

	router := gin.Default()

	// Endpoints
	router.GET("/items", func(c *gin.Context) {
		items, err := service.GetAllItems()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	router.GET("/items/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, err := service.GetItemByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	router.POST("/items", func(c *gin.Context) {
		var item domain.Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		service.CreateItem(item)
		c.JSON(http.StatusCreated, item)
	})

	router.PUT("/items/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedItem domain.Item
		if err := c.ShouldBindJSON(&updatedItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedItem.ID = id
		err := service.UpdateItem(updatedItem)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedItem)
	})

	router.DELETE("/items/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := service.DeleteItem(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
	})

	router.Run(":8000")
}
