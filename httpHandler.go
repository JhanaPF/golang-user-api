package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s %s", dictionaryApiUrl, endpoint)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête :", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return nil, err
	}

	fmt.Println("Réponse :", string(body)) // Convertir le corps de réponse en une chaîne de caractères et l'afficher
	return body, err
}
