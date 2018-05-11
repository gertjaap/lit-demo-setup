package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-docker-tester/btc"
)

func MineBlockHandler(w http.ResponseWriter, r *http.Request) {

	err := btc.MineBlocks(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
