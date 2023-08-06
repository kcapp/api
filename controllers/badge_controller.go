package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kcapp/api/data"
)

func GetBadges(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	badges, err := data.GetBadges()
	if err != nil {
		log.Println("Unable to get badges")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(badges)
}
