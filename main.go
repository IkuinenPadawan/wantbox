package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WishListItem struct {
	ID       string  `json:"id"`
	ItemName string  `json:"item"`
	Url      string  `json:"url"`
	Price    float64 `json:"price"`
}

type WishListPageData struct {
	PageTitle     string
	WishlistItems []WishListItem
}

var wishlist = []WishListItem{
	{ID: "1", ItemName: "Rudder Pedals", Url: "www.google.com", Price: 350.00},
	{ID: "2", ItemName: "Synth", Url: "www.google.com", Price: 500.00},
	{ID: "3", ItemName: "Saucony shoes", Url: "www.google.com", Price: 99.00},
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("layout.html")

	router.GET("/", func(c *gin.Context) {
		data := WishListPageData{
			PageTitle:     "Wantbox",
			WishlistItems: wishlist,
		}
		c.HTML(http.StatusOK, "layout.html", data)
	})

	router.Run(":8080")
}
