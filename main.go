package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Go Postgres driver
)

var dictionaryApiUrl string = "http://localhost:7005/"
var db *sql.DB

func getUserIdFromToken(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Remplacez "secret" par votre clé secrète
	})

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))
	return userID
}

func init() {
	db, err := sql.Open("postgres", "user=root password=password dbname=user-golang-api host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)

	// Routes pour les fonctionnalités spécifiques de l'utilisateur
	authorized := router.Group("/api")
	authorized.Use(authMiddleware)
	{
		authorized.GET("/history", getHistory)
		authorized.GET("/friends", getFriends)
	}

	router.Run("localhost:8080")
	port := os.Getenv("PORT")
	router.Run(fmt.Sprintf("%s %s", ":", port))
}
