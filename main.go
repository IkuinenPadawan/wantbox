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

type WishlistItemEditPageData struct {
	PageTitle    string
	WishlistItem WishListItem
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

func deleteWishlistItemHandler(db *sql.DB, c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
	}
	deleteWishlistItem(db, id)

	c.Redirect(http.StatusFound, "/")
}

func deleteWishlistItem(db *sql.DB, itemId int) {
	deleteOneSQL := `DELETE FROM wishlist WHERE ID = (?)`
	statement, err := db.Prepare(deleteOneSQL)
	if err != nil {
		log.Print(err)
	}
	_, err = statement.Exec(itemId)
	if err != nil {
		log.Print(err)
	}
}

func updateWishlistItemHandler(db *sql.DB, c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
	}
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

	updateWishlistItem(db, itemname, price, url, user, id)
	c.Redirect(http.StatusFound, "/")
}

func updateWishlistItem(db *sql.DB, itemname string, price float64, url string, user int, id int) {
	updateSQL := `UPDATE wishlist SET itemname = (?), price = (?), url = (?), user_id = (?) WHERE ID = (?)`
	statement, err := db.Prepare(updateSQL)

	if err != nil {
		log.Print(err)
	}
	_, err = statement.Exec(itemname, price, url, user, id)
	if err != nil {
		log.Print(err)
	}
}

func findAll(db *sql.DB) ([]WishListItem, error) {
	selectAll := `SELECT id, itemname, price, url FROM wishlist;`
	rows, err := db.Query(selectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlistItems []WishListItem
	for rows.Next() {
		w := &WishListItem{}
		err := rows.Scan(&w.ID, &w.ItemName, &w.Price, &w.Url)
		if err != nil {
			return nil, err
		}
		wishlistItems = append(wishlistItems, *w)
	}
	return wishlistItems, nil
}

func findOne(db *sql.DB, id int) (WishListItem, error) {
	selectOne := `SELECT id, itemname, price, url, user_id FROM wishlist WHERE ID = ?`
	statement, err := db.Prepare(selectOne)
	if err != nil {
		log.Print(err)
	}
	defer statement.Close()

	var wishlistItem WishListItem
	w := &WishListItem{}
	err = statement.QueryRow(id).Scan(&w.ID, &w.ItemName, &w.Price, &w.Url, &w.UserId)
	if err != nil {
		log.Print(err)
	}
	wishlistItem = *w
	return wishlistItem, nil

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
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		wishlist, err := findAll(db)
		if err != nil {
			log.Fatal("Error finding all items:", err)
		}
		data := WishListPageData{
			PageTitle:     "Wantbox",
			WishlistItems: wishlist,
		}
		c.HTML(http.StatusOK, "layout.tmpl", data)
	})

	router.POST("/wishlist", func(c *gin.Context) {
		handleWishlistItemForm(db, c)
	})

	router.POST("/wishlist/:id/delete", func(c *gin.Context) {
		deleteWishlistItemHandler(db, c)
	})

	router.GET("/wishlist/:id/edit", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		wishlistItem, err := findOne(db, id)
		if err != nil {
			log.Print(err)
		}
		data := WishlistItemEditPageData{
			PageTitle:    "Edit",
			WishlistItem: wishlistItem,
		}
		c.HTML(http.StatusOK, "edititem.tmpl", data)
	})

	router.POST("/wishlist/:id/edit", func(c *gin.Context) {
		updateWishlistItemHandler(db, c)
	})

	router.Run(":8080")
}
