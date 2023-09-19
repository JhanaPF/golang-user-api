package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getFriends(c *gin.Context) {
	userID := getUserIdFromToken(c)

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
