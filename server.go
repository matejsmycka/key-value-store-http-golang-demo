package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB(dbConnStr string) {
	if dbConnStr == "" {
		log.Fatal("DB_CONN_STR not set")
	}

	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        key TEXT NOT NULL,
        value TEXT NOT NULL
    ); `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

type Item struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func saveItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	query := `INSERT INTO items (key, value) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, item.Key, item.Value).Scan(&item.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert item"})
		return
	}

	c.JSON(201, item)
}

func getItems(c *gin.Context) {
	// returns [] if no items found
	rows, err := db.Query("SELECT id, key, value FROM items")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Key, &item.Value); err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch items"})
			return
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}

func getIndex(c *gin.Context) {
	c.File("index.html")
}

func deleteItem(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		c.JSON(400, gin.H{"error": "Key is required"})
		return
	}

	_, err := db.Exec("DELETE FROM items WHERE key = $1", key)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(200, gin.H{"message": "Item deleted successfully"})
}

func main() {
	port := "8080"
	initDB(os.Getenv("DB_CONN_STR"))
	log.Println("DB connection initialized")
	defer db.Close()

	r := gin.Default()
	r.POST("/items", saveItem)
	r.GET("/items", getItems)
	r.GET("/", getIndex)
	r.DELETE("/items", deleteItem)

	log.Println("Server running on port", port)
	r.Run(":" + port)
}
