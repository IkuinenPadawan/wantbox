package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type wishlisting struct {
	ID    string  `json:"id"`
	Item  string  `json:"item"`
	Url   string  `json:"url"`
	Price float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, wishlist)
}

var wishlist = []wishlisting{
	{ID: "1", Item: "Rudder Pedals", Url: "www.google.com", Price: 350.00},
	{ID: "2", Item: "Synth", Url: "www.google.com", Price: 500.00},
	{ID: "3", Item: "Saucony shoes", Url: "www.google.com", Price: 99.00},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8080")
}
