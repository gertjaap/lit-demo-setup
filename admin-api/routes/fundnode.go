package routes

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := commands.GetAddresses(rpcCon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, adr := range addr.WitAddresses {
		logging.Info.Printf("Found address: %s\n", adr)
	}
	/*
		// Send it a bunch each coin
		for _, cd := range coindaemon.CoinDaemons {
			rpc, err := coindaemon.GetRpcClient(cd)
			rpc.send

		}

		logging.Info.Printf("Funding [%s] with 1 BTC-r", addr)
		err = btc.SendCoins(addr.WitAddresses, 100000000) // 1 BTC-r
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}*/

	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
