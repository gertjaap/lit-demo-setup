package routes

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"

	"github.com/gertjaap/lit-demo-setup/admin-api/models"
)

func ListNodesHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		logging.Error.Printf("ListNodesHandler NewEnvClient error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	containers, err := docker.LitNodes(cli)
	if err != nil {
		logging.Error.Printf("ListNodesHandler LitNodes error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ditch the big fat node from the list
	filteredContainers := make([]types.Container, 0)
	for _, c := range containers {
		if c.Names[0][1:] != "litdemobigfatnode" {
			filteredContainers = append(filteredContainers, c)
		}
	}

	// Connect to all nodes and fetch their status
	nodes := make([]models.LitNode, len(filteredContainers))

	for i, c := range filteredContainers {
		nodes[i].Name = c.Names[0][1:]
		rpcCon, err := docker.GetLndcRpc(cli, nodes[i].Name, false)
		if err != nil {
			logging.Error.Printf("ListNodesHandler GetLndcRpc error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nodes[i].Balances, err = litrpc.GetBalancesFromNode(rpcCon)
		if err != nil {
			logging.Error.Printf("ListNodesHandler GetBalancesFromNode error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := docker.NodeAddresses[nodes[i].Name]; !ok {
			docker.NodeAddresses[nodes[i].Name], err = docker.GetAddress(cli, nodes[i].Name)
			if err != nil {
				logging.Error.Printf("ListNodesHandler GetAddress error: %s", err.Error())
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
		logging.Error.Printf("ListNodesHandler json.Marshal error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
