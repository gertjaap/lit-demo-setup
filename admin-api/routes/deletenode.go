package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gorilla/mux"
)

func DeleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	vars := mux.Vars(r)
	err = docker.DropLitNode(cli, vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dataDir := fmt.Sprintf("/data/%s", vars["id"][7:])
	logging.Info.Printf("Dropping datadir %s", dataDir)
	err = os.RemoveAll(dataDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodesLastRefreshed = time.Now().Add(-30 * time.Second) // invalidate cache

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
