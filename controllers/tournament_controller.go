package controllers

import (
	"encoding/json"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
)

// GetTournaments will return all tournaments
func GetTournaments(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	tournaments, err := data.GetTournaments()
	if err != nil {
		log.Println("Unable to get tournaments", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournaments)
}
