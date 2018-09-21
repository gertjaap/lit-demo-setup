package routes

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"

	"github.com/gertjaap/lit-demo-setup/admin-api/models"
)

func ListNodesHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	containers, err := docker.LitNodes(cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Connect to all nodes and fetch their status

	nodes := make([]models.LitNode, len(containers))

	for i, c := range containers {
		nodes[i].Name = c.Names[0][1:]
		rpcCon, err := docker.GetLndcRpc(cli, nodes[i].Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nodes[i].Balances, err = litrpc.GetBalancesFromNode(rpcCon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := docker.NodeAddresses[nodes[i].Name]; !ok {
			docker.NodeAddresses[nodes[i].Name], err = docker.GetAddress(cli, nodes[i].Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		nodes[i].Address = docker.NodeAddresses[nodes[i].Name]

		for _, p := range c.Ports {
			if p.PrivatePort == 2448 {
				nodes[i].PublicLitPort = int(p.PublicPort)
			}
		}
	}

	js, err := json.Marshal(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
