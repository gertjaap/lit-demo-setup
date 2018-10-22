package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/models"
)

var nodesLastRefreshed time.Time
var cachedNodes []models.LitNode

func GetCachedNodes() []models.LitNode {
	return cachedNodes
}

func CacheNodes() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	defer cli.Close()
	containers, err := docker.LitNodes(cli)
	if err != nil {
		return err
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
			return err
		}
		nodes[i].Balances, err = litrpc.GetRawBalancesFromNode(rpcCon)
		if err != nil {
			logging.Error.Printf("Error while refreshing balance of %s: %s\n", nodes[i].Name, err.Error())
			nodes[i].Error = true
			nodes[i].ErrorDetails = err.Error()
			if strings.Contains(err.Error(), "timeout") {
				// If we receive an RPC timeout, just try it in the next run.
				continue
			}
			// If something goes wrong on the rpc connection, reset it.
			err = rpcCon.Reconnect()
			if err != nil {
				// remove the lndc
				docker.DropLndcRpc(nodes[i].Name)
			}

			continue
		}

		if _, ok := docker.NodeAddresses[nodes[i].Name]; !ok {
			docker.NodeAddresses[nodes[i].Name], err = docker.GetAddress(cli, nodes[i].Name)
			if err != nil {
				nodes[i].Error = true
				nodes[i].ErrorDetails = err.Error()
				continue
			}
		}
		nodes[i].Address = docker.NodeAddresses[nodes[i].Name]
		nodes[i].TrackerOK = true
		nodes[i].TrackerIP, err = trackerLookup(nodes[i].Address)
		if err != nil {
			nodes[i].Error = true
			nodes[i].ErrorDetails = err.Error()
			nodes[i].TrackerOK = false
			continue
		}

		for _, p := range c.Ports {
			if p.PrivatePort == 2448 {
				nodes[i].PublicLitPort = int(p.PublicPort)
			}
		}
	}
	cachedNodes = nodes
	nodesLastRefreshed = time.Now()

	return nil
}

func ListNodesHandler(w http.ResponseWriter, r *http.Request) {

	js, err := json.Marshal(cachedNodes)
	if err != nil {
		logging.Error.Printf("ListNodesHandler json.Marshal error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type nodeinfo struct {
	Success bool
	Node    struct {
		IPv4 string
		IPv6 string
		Addr string
	}
}

func trackerLookup(adr string) (string, error) {
	var client http.Client

	resp, err := client.Get(fmt.Sprintf("http://litdemotracker:46580/%s", adr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var node nodeinfo
	err = decoder.Decode(&node)
	if err != nil {
		return "", err
	}

	if !node.Success {
		return "", fmt.Errorf("Node not found")
	}

	return node.Node.IPv4, nil

}
