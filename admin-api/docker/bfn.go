package docker

import (
	"context"
	"strings"
	"sync"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemon"
	"github.com/mit-dci/lit/litrpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
)

var bfnConnectMutex sync.Mutex

func InitBigFatNode(cli *client.Client) error {
	bfnConnectMutex = sync.Mutex{}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	logging.Info.Printf("Checking if big fat node is running\n")
	found := false
	for _, c := range containers {
		correctNetwork := false
		for _, n := range c.NetworkSettings.Networks {
			if n.NetworkID == NetworkID {
				correctNetwork = true
			}
		}
		if !correctNetwork {
			continue
		}
		if c.Names[0][1:] == "litdemobigfatnode" {
			found = true
		}
	}

	if !found {
		logging.Info.Printf("Big fat node not found, creating...\n")
		node, err := NewLitNode(cli)
		if err != nil {
			return err
		}

		err = cli.ContainerRename(context.Background(), node.ID, "litdemobigfatnode")
		if err != nil {
			return err
		}

		// Fund the big fat node
		rpcClient, err := GetLndcRpc(cli, "litdemobigfatnode")
		if err != nil {
			return err
		}

		addrs, err := commands.GetAddresses(rpcClient)
		if err != nil {
			return err
		}

		logging.Info.Println("Funding Big fat node on these addresses:")

		for _, adr := range addrs.WitAddresses {
			logging.Info.Printf("Funding on %s\n", adr)
			for _, cd := range coindaemon.CoinDaemons {
				if strings.HasPrefix(adr, cd.CoinParams.Bech32Prefix) {
					// Send 50 transactions of 100 coins of each
					for i := 0; i < 50; i++ {
						err = cd.SendCoins(adr, 10000000000)
						if err != nil {
							return err
						}
					}

				}
			}
		}
	}
	return nil
}

func ConnectBFNToNode(rpcClient *litrpc.LndcRpcClient, node string) (uint32, error) {
	bfnConnectMutex.Lock()

	conns1, err := commands.ListConnections(rpcClient)
	if err != nil {
		return 0, err
	}

	_, err = commands.Connect(rpcClient, NodeAddresses[node])
	if err != nil {
		return 0, err
	}

	conns2, err := commands.ListConnections(rpcClient)
	if err != nil {
		return 0, err
	}

	returnVal := uint32(0)

	for _, c2 := range conns2.Connections {
		found := false
		for _, c1 := range conns1.Connections {
			if c1.PeerNumber == c2.PeerNumber {
				found = true
			}
		}
		if !found {
			returnVal = c2.PeerNumber
		}
	}

	bfnConnectMutex.Unlock()

	return returnVal, nil
}
