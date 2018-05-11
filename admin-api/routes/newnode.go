package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gertjaap/blockchain-indexer-insight/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
)

func NewNodeHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	node, err := docker.NewLitNode(cli)
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rpcUrl := fmt.Sprintf("ws://%s:8001/ws", node.Name)
	wsConn, rpcCon, err := litrpc.ConnectWithRetry(rpcUrl)
	defer wsConn.Close()
	node.Address, err = litrpc.GetAddressFromNode(rpcCon)
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	docker.NodeAddresses[node.Name] = node.Address
	node.Balances = map[uint32]int64{257: int64(0)}
	js, err := json.Marshal(node)
	if err != nil {
		logging.Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
