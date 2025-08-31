package main

import (
	"database/sql"
	"fmt"
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
	UserName string  `json:"username"`
}

type WishListPageData struct {
	PageTitle     string
	WishlistItems []WishListItem
	Users         []User
}

type WishlistItemEditPageData struct {
	PageTitle    string
	WishlistItem WishListItem
	Users        []User
}

type AddUserPageData struct {
	PageTitle string
}

type ValidatedWishListItem struct {
	ItemName string
	Price    float64
	Url      string
	UserID   int
}

type ValidatedUserName struct {
	UserName string
}

type User struct {
	ID   int
	Name string
}

func validateAndParseWishlistitem(itemname string, priceStr string, url string, userStr string) (ValidatedWishListItem, error) {
	var item ValidatedWishListItem
	if itemname == "" || len(itemname) > 100 {
		return item, fmt.Errorf("item name needs to be between 1-100 characters")
	}
	item.ItemName = itemname

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return item, fmt.Errorf("invalid price")
	}
	item.Price = price

	if len(url) > 500 {
		return item, fmt.Errorf("url needs to be under 500 characters")
	}
	item.Url = url

	userId, err := strconv.Atoi(userStr)
	if err != nil {
		return item, fmt.Errorf("invalid user id")
	}
	item.UserID = userId

	return item, nil
}

func validateUserForm(username string) (ValidatedUserName, error) {
	var validatedUsername ValidatedUserName
	if username == "" || len(username) > 40 {
		return validatedUsername, fmt.Errorf("username needs to be between 1-40 charactes")
	}
	validatedUsername.UserName = username
	return validatedUsername, nil
}

func handleWishlistItemForm(db *sql.DB, c *gin.Context) {
	itemname := c.PostForm("itemname")
	priceStr := c.PostForm("price")
	url := c.PostForm("url")
	userStr := c.PostForm("user")

	validatedItem, err := validateAndParseWishlistitem(itemname, priceStr, url, userStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = insertWishlistItem(db, validatedItem.ItemName, validatedItem.Price, validatedItem.Url, validatedItem.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save item"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func insertWishlistItem(db *sql.DB, itemname string, price float64, url string, user int) error {
	insertWishlistItemSQL := `INSERT INTO wishlist (itemname, price, url, user_id) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertWishlistItemSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(itemname, price, url, user)

	return err
}

func deleteWishlistItemHandler(db *sql.DB, c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id format"})
		return
	}
	err = deleteWishlistItem(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func deleteWishlistItem(db *sql.DB, itemId int) error {
	deleteOneSQL := `DELETE FROM wishlist WHERE ID = (?)`
	statement, err := db.Prepare(deleteOneSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(itemId)
	return err
}

func updateWishlistItemHandler(db *sql.DB, c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id format"})
		return
	}
	itemname := c.PostForm("itemname")
	priceStr := c.PostForm("price")
	url := c.PostForm("url")
	userStr := c.PostForm("user")

	validatedItem, err := validateAndParseWishlistitem(itemname, priceStr, url, userStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateWishlistItem(db, validatedItem.ItemName, validatedItem.Price, validatedItem.Url, validatedItem.UserID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func updateWishlistItem(db *sql.DB, itemname string, price float64, url string, user int, id int) error {
	updateSQL := `UPDATE wishlist SET itemname = (?), price = (?), url = (?), user_id = (?) WHERE ID = (?)`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(itemname, price, url, user, id)

	return err
}

func handleUserForm(db *sql.DB, c *gin.Context) {
	username := c.PostForm("username")
	validatedUserName, err := validateUserForm(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = addUser(db, validatedUserName.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func addUser(db *sql.DB, username string) error {
	addUserSQL := `INSERT INTO users(name) VALUES (?);`
	statement, err := db.Prepare(addUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(username)

	return err
}

func findAllUsers(db *sql.DB) ([]User, error) {
	selectAll := `SELECT id, name FROM users;`
	rows, err := db.Query(selectAll)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		users = append(users, *u)
	}

	return users, nil
}

func findAll(db *sql.DB) ([]WishListItem, error) {
	selectAll := `SELECT wishlist.id, itemname, price, url, users.name FROM wishlist INNER JOIN users ON users.id = wishlist.user_id;`
	rows, err := db.Query(selectAll)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	var wishlistItems []WishListItem
	for rows.Next() {
		w := &WishListItem{}
		err := rows.Scan(&w.ID, &w.ItemName, &w.Price, &w.Url, &w.UserName)
		if err != nil {
			log.Print(err)
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
	wishlistItem = *w
	return wishlistItem, err
}

func main() {
	// DB
	db, err := sql.Open("sqlite", "wantbox.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dropWishlist := `DROP TABLE wishlist;`
	dropUsers := `DROP TABLE users;`
	_, err = db.Exec(dropWishlist)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(dropUsers)
	if err != nil {
		log.Fatal(err)
	}

	createWishlistTableSql := `CREATE TABLE IF NOT EXISTS wishlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemname TEXT NOT NULL,
		price REAL NOT NULL,
		url TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	    FOREIGN KEY (user_id)
	      REFERENCES users (id)
		);`

	insertMockData := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Rudder Pedals", 350.00, "www.google.com", 1);`
	insertMockData2 := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Synth", 500.00, "www.google.com", 1);`
	insertMockData3 := `INSERT INTO wishlist(itemname, price, url, user_id) VALUES("Saucony shoes", 99.00, "www.google.com", 2);`

	createUserTableSql := `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL);`

	insertMockUserData := `INSERT INTO users(name) VALUES ("Augustus");`
	insertMockUserData2 := `INSERT INTO users(name) VALUES ("Magalhaes");`

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

	_, err = db.Exec(createUserTableSql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(insertMockUserData)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(insertMockUserData2)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
			return
		}

		users, err := findAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		data := WishListPageData{
			PageTitle:     "Wantbox",
			WishlistItems: wishlist,
			Users:         users,
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
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		wishlistItem, err := findOne(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve item"})
			}
			return
		}

		users, err := findAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		data := WishlistItemEditPageData{
			PageTitle:    "Wantbox",
			WishlistItem: wishlistItem,
			Users:        users,
		}
		c.HTML(http.StatusOK, "edititem.tmpl", data)
	})

	router.POST("/wishlist/:id/edit", func(c *gin.Context) {
		updateWishlistItemHandler(db, c)
	})

	router.GET("/user", func(c *gin.Context) {
		data := AddUserPageData{
			PageTitle: "Wantbox",
		}
		c.HTML(http.StatusOK, "adduser.tmpl", data)
	})

	router.POST("/user", func(c *gin.Context) {
		handleUserForm(db, c)
	})

	router.Run(":8089")
}
