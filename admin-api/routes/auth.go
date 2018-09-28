package routes

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gorilla/mux"
)

func AuthNodeHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	vars := mux.Vars(r)
	rpcCon, err := docker.GetLndcRpc(cli, vars["id"])
	if err != nil {
		logging.Error.Printf("AuthNodeHandler GetLndcRpc error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubKey, err := hex.DecodeString(vars["key"])
	if err != nil {
		logging.Error.Printf("AuthNodeHandler hex.DecodeString error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pubKey33 [33]byte
	copy(pubKey33[:], pubKey)

	err = litrpc.RcAuth(rpcCon, pubKey33, vars["yesno"] == "1")
	if err != nil {
		logging.Error.Printf("AuthNodeHandler RcAuth error: %s", err.Error())
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

func PendingAuthForNodeHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	vars := mux.Vars(r)
	rpcCon, err := docker.GetLndcRpc(cli, vars["id"])
	if err != nil {
		logging.Error.Printf("PendingAuthForNodeHandler GetLndcRpc error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubKeys, err := litrpc.PendingRCAuthRequests(rpcCon)
	if err != nil {
		logging.Error.Printf("PendingAuthForNodeHandler PendingRCAuthRequests error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]string, len(pubKeys))
	for i, k := range pubKeys {
		result[i] = hex.EncodeToString(k[:])
	}
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
