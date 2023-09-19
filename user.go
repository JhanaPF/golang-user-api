package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Go Postgres driver
	"golang.org/x/crypto/bcrypt"
)

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token manquant"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Remplacez "secret" par une clé secrète plus sécurisée
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token non valide"})
		c.Abort()
		return
	}

	c.Next()
}

func register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Données non valides"})
		return
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de hachage du mot de passe"})
		return
	}

	// Enregistrement de l'utilisateur dans la base de données
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur enregistré avec succès"})
}

func login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Données non valides"})
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", input.Username).Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Nom d'utilisateur ou mot de passe incorrect"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Nom d'utilisateur ou mot de passe incorrect"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = input.Username
	tokenString, err := token.SignedString([]byte("secret")) // Remplacez "secret" par une clé secrète plus sécurisée
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
