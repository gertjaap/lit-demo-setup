package routes

import (
	"encoding/json"
	"fmt"
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
		rpcUrl := fmt.Sprintf("ws://%s:8001/ws", nodes[i].Name)
		wsConn, rpcCon, err := litrpc.Connect(rpcUrl)
		nodes[i].Balances, err = litrpc.GetBalancesFromNode(rpcCon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if val, ok := docker.NodeAddresses[nodes[i].Name]; ok {
			nodes[i].Address = val
		} else {
			if err == nil {
				defer wsConn.Close()
				nodes[i].Address, err = litrpc.GetAddressFromNode(rpcCon)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				docker.NodeAddresses[nodes[i].Name] = nodes[i].Address
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		for _, p := range c.Ports {
			if p.PrivatePort == 2448 {
				nodes[i].PublicLitPort = int(p.PublicPort)
			}
			if p.PrivatePort == 8001 {
				nodes[i].PublicRpcPort = int(p.PublicPort)
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
