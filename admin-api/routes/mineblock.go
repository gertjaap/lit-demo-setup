package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
)

func MineBlockHandler(w http.ResponseWriter, r *http.Request) {
	err := coindaemons.MineBlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
