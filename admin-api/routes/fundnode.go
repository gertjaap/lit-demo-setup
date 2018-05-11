package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-docker-tester/btc"
	"github.com/gorilla/mux"
)

func FundNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rpcUrl := fmt.Sprintf("ws://%s:8001/ws", vars["id"])
	wsConn, rpcCon, err := litrpc.Connect(rpcUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer wsConn.Close()
	addr, err := litrpc.GetBlockchainAddress(rpcCon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info.Printf("Funding [%s] with 1 BTC-r", addr)
	err = btc.SendCoins(addr, 100000000) // 1 BTC-r
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
