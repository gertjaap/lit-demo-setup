package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
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
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info.Printf("Funding [%s] with 1 BTC-r", addr)

	client, err := btc.GetRpcClient()
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	btcAddr, err := btcutil.DecodeAddress(addr, &chaincfg.RegressionNetParams)
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txHash, err := client.SendFrom("", btcAddr, btcutil.Amount(100000000))
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Info.Println("Sent coins, tx hash:", txHash)

	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
