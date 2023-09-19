package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Go Postgres driver
)

func getHistory(c *gin.Context) {
	userID := getUserIdFromToken(c)

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
