package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemon"
)

func MineBlockHandler(w http.ResponseWriter, r *http.Request) {

	for _, cd := range coindaemon.CoinDaemons {
		err := cd.MineBlocks(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
