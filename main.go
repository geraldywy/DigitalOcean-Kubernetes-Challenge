package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	initDB()
	router := gin.Default()

	router.GET("/messages", getMessages)
	router.POST("/message", addMessage)

	router.Run(":3000")
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:password123@tcp(localhost:6446)/example_db") // hardcoded
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func addMessage(c *gin.Context) {
	msg := new(Message)
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("adding msg: %+v\n", msg)

	_, err := db.ExecContext(c, "INSERT INTO messages VALUES (?, ?)", msg.Author, msg.Message)
	if err != nil {
		log.Printf("err inserting msg: %+v, err: %v\n", msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": msg})
}

func getMessages(c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT * FROM messages")
	if err != nil {
		log.Printf("err fetching messages from db: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	msgs := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		if err := rows.Scan(&msg.Author, &msg.Message); err != nil {
			log.Printf("err reading message, err: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		msgs = append(msgs, msg)
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

type Message struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}
