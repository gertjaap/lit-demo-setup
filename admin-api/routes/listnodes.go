package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func ListNodesHandler(w http.ResponseWriter, r *http.Request) {
	if nodesLastRefreshed.Add(20 * time.Second).Before(time.Now()) { // Cache expired.
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
			nodes[i].Balances, err = litrpc.GetRawBalancesFromNode(rpcCon)
			if err != nil {
				nodes[i].Error = true
				nodes[i].ErrorDetails = err.Error()

				// If something goes wrong on the rpc connection, reset it.
				err = rpcCon.Reconnect()
				if err != nil {
					// remove the lndc
					docker.DropLndcRpc(cli, nodes[i].Name)
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

			nodes[i].Channels, err = litrpc.GetChannelsFromNode(rpcCon)
			if err != nil {
				nodes[i].Error = true
				nodes[i].ErrorDetails = err.Error()
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
	}

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
