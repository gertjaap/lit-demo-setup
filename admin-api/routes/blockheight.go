package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-docker-tester/btc"
)

func BlockHeightHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := btc.GetRpcClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blocks, err := cli.GetBlockCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
