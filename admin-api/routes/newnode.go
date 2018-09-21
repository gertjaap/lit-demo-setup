package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemon"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
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

	rpcClient, err := docker.GetLndcRpc(cli, "litdemobigfatnode")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Connect the BFN to the new node
	peerIdx, err := docker.ConnectBFNToNode(rpcClient, node.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, cd := range coindaemon.CoinDaemons {
		reply, err := commands.Fund(rpcClient, peerIdx, cd.LitCoinType, 4000000000, 500000000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logging.Info.Printf("Funded %s - %s\n", cd.ContainerName, reply.Status)
	}

	js, err := json.Marshal(node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
