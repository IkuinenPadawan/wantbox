package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
)

type WishListItem struct {
	ID       string  `json:"id"`
	ItemName string  `json:"itemname"`
	Url      string  `json:"url"`
	Price    float64 `json:"price"`
	UserId   int     `json:"user_id"`
}

type WishListPageData struct {
	PageTitle     string
	WishlistItems []WishListItem
}

func handleWishlistItemForm(db *sql.DB, c *gin.Context) {
	itemname := c.PostForm("itemname")
	priceStr := c.PostForm("price")
	url := c.PostForm("url")
	userStr := c.PostForm("user")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Print(err)
	}

	user, err := strconv.Atoi(userStr)
	if err != nil {
		log.Print(err)
	}

	insertWishlistItem(db, itemname, price, url, user)
	c.Redirect(http.StatusFound, "/")
}

func insertWishlistItem(db *sql.DB, username string, price float64, url string, user int) {
	insertWishlistItemSQL := `INSERT INTO wishlist (itemname, price, url, user_id) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertWishlistItemSQL)
	if err != nil {
		log.Print(err)
	}
	_, err = statement.Exec(username, price, url, user)
	if err != nil {
		log.Print(err)
	}
}

func findAll(db *sql.DB) ([]WishListItem, error) {
	selectAll := `SELECT itemname, price, url FROM wishlist;`
	rows, err := db.Query(selectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlistItems []WishListItem
	for rows.Next() {
		w := &WishListItem{}
		err := rows.Scan(&w.ItemName, &w.Price, &w.Url)
		if err != nil {
			return nil, err
		}
		wishlistItems = append(wishlistItems, *w)
	}
	return wishlistItems, nil
}

func main() {
	// DB
	db, err := sql.Open("sqlite", "wantbox.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createWishlistTableSql := `CREATE TABLE IF NOT EXISTS wishlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemname TEXT NOT NULL,
		price REAL NOT NULL,
		url TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

	insertMockData := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Rudder Pedals", 350.00, "www.google.com", 1);`
	insertMockData2 := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Synth", 500.00, "www.google.com", 1);`
	insertMockData3 := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Saucony shoes", 99.00, "www.google.com", 1);`

	_, err = db.Exec(createWishlistTableSql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(insertMockData)
	_, err = db.Exec(insertMockData2)
	_, err = db.Exec(insertMockData3)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Router
	router := gin.Default()
	router.LoadHTMLFiles("layout.html")

	router.GET("/", func(c *gin.Context) {
		wishlist, err := findAll(db)
		if err != nil {
			log.Fatal("Error finding all items:", err)
		}
		data := WishListPageData{
			PageTitle:     "Wantbox",
			WishlistItems: wishlist,
		}
		c.HTML(http.StatusOK, "layout.html", data)
	})

	router.POST("/wishlist", func(c *gin.Context) {
		handleWishlistItemForm(db, c)
	})

	router.Run(":8080")
}
