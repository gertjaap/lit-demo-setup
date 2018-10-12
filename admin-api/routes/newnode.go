package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
)

func NewNodeHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	node, err := docker.NewLitNode(cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	node.Address, err = docker.GetAddress(cli, node.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	docker.NodeAddresses[node.Name] = node.Address

	err = docker.ConnectAndFund(cli, node.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodesLastRefreshed = time.Now().Add(-30 * time.Second) // invalidate cache

	js, err := json.Marshal(node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
