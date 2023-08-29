package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "A song of ice and fire", Author: "George RR Martin", Quantity: 1000},
	{ID: "2", Title: "The lord of the rings", Author: "J RR Tolkien", Quantity: 2000},
}

func getBooks(context *gin.Context) {
	rows, err := db.Query("SELECT * FROM yourtable")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Lire les données résultantes ici
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	context.IndentedJSON(http.StatusOK, books)
}

func createBook(context *gin.Context) {
	var newBook book
	if err := context.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

var dictionaryApiKey string = "http://localhost:7005/"

func request(endpoint string) ([]byte, error) {
	// Effectuer une requête GET à l'URL spécifiée
	url := fmt.Sprintf("%s %s", dictionaryApiKey, endpoint)
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

func getCourses(context *gin.Context) {

}

var db *sql.DB

func init() {
	db, err := sql.Open("postgres", "user=root password=password dbname=test-golang-api host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func main() {
	router := gin.Default()
	router.GET("/courses", getCourses)
	//router.GET("/lessons", getLessons)
	//router.GET("/questions", getQuestions)
	router.Run("localhost:8080")
	port := os.Getenv("PORT")
	router.Run(fmt.Sprintf("%s %s", ":", port))
}
