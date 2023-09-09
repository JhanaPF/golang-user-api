package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Go Postgres driver
	"golang.org/x/crypto/bcrypt"
)

var dictionaryApiUrl string = "http://localhost:7005/"

func request(endpoint string) ([]byte, error) {
	// Effectuer une requête GET à l'URL spécifiée
	url := fmt.Sprintf("%s %s", dictionaryApiUrl, endpoint)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête :", err)
		return nil, err
	}
	defer response.Body.Close()

	// Lire le contenu de la réponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return nil, err
	}

	// Convertir le corps de réponse en une chaîne de caractères et l'afficher
	fmt.Println("Réponse :", string(body))
	return body, err
}

func getHistory(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur à partir du token
	userID := getUserIdFromToken(c)

	// Requête à la base de données pour récupérer l'historique des parties de l'utilisateur
	rows, err := db.Query("SELECT id, user_id FROM game_history WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la récupération de l'historique des parties"})
		return
	}
	defer rows.Close()

	var history []GameHistory
	for rows.Next() {
		var game GameHistory
		if err := rows.Scan(&game.ID, &game.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la lecture des données de l'historique des parties"})
			return
		}
		history = append(history, game)
	}

	c.JSON(http.StatusOK, history)
}

func getFriends(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur à partir du token
	userID := getUserIdFromToken(c)

	// Requête à la base de données pour récupérer la liste d'amis de l'utilisateur
	rows, err := db.Query("SELECT id, user_id, friend_id FROM dictionaryApiUrls WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la récupération de la liste d'amis"})
		return
	}
	defer rows.Close()

	var friends []Relationship
	for rows.Next() {
		var friend Relationship
		if err := rows.Scan(&friend.ID, &friend.UserID, &friend.FriendID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la lecture des données de la liste d'amis"})
			return
		}
		friends = append(friends, friend)
	}

	c.JSON(http.StatusOK, friends)
}

// Fonction utilitaire pour récupérer l'ID de l'utilisateur à partir du token JWT
func getUserIdFromToken(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Remplacez "secret" par votre clé secrète
	})

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64)) // L'ID de l'utilisateur est stocké dans le token
	return userID
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

	// Créer un jeton JWT
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

var db *sql.DB

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

	// Créer les tables si elles n'existent pas
	//db.AutoMigrate(&User{}, &dictionaryApiUrl{}, &GameHistory{}, &CourseChoice{})

	//_, err = db.Exec(`
	//   CREATE TABLE IF NOT EXISTS utilisateurs (
	//	   id SERIAL PRIMARY KEY,
	//	   nom VARCHAR(255),
	//	   prenom VARCHAR(255)
	//   )
	//`)
	//if err != nil {
	//   log.Fatal(err)
	//}
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
