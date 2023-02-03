package controllers

import (
	"encoding/json"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
)

// GetOweTypes will return all owe types
func GetOweTypes(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	owes, err := data.GetOweTypes()
	if err != nil {
		log.Println("Unable to get owe types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(owes)
}
