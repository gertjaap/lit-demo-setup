package routes

import (
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gorilla/mux"
)

func FundNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	rpcCon, err := docker.GetLndcRpc(cli, vars["id"])
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info.Printf("Funding [%s] with 5 BTC-r", addr)

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

	txHash, err := client.SendFrom("", btcAddr, btcutil.Amount(500000000))
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
