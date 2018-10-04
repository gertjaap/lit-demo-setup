package docker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gertjaap/lit-demo-setup/admin-api/litrpc"
	litrpclit "github.com/mit-dci/lit/litrpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
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
		rpcClient, err := GetLndcRpc(cli, "litdemobigfatnode", true)
		if err != nil {
			return err
		}

		// Authorize admin panel on bigfatnode
		adminPanelKey, err := GetAdminPanelKey()
		if err != nil {
			logging.Error.Printf("Error getting admin panel key: %s\n", err.Error())
			return err
		}
		adminPanelPubKey := [33]byte{}
		copy(adminPanelPubKey[:], adminPanelKey.PubKey().SerializeCompressed())

		err = litrpc.RcAuth(rpcClient, adminPanelPubKey, true)
		if err != nil {
			logging.Error.Printf("ConnectAndFund RcAuth error: %s", err.Error())
			return err
		}

		addrs, err := commands.GetAddresses(rpcClient)
		if err != nil {
			return err
		}

		logging.Info.Println("Funding Big fat node on these addresses:")

		for _, adr := range addrs.WitAddresses {
			logging.Info.Printf("Funding on %s\n", adr)
			for _, cd := range coindaemons.CoinDaemons {
				if strings.HasPrefix(adr, cd.CoinParams.Bech32Prefix) {
					// Send 100 transactions of InitialFunding/100 coins of each
					for i := 0; i < 100; i++ {
						err = cd.SendCoins(adr, int64(cd.InitialFunding)*int64(1000000))
						if err != nil {
							return err
						}
					}

				}
			}
		}

		// Mine a few blocks to confirm new funds
		for i := 0; i < 10; i++ {
			coindaemons.MineBlock()
		}

		rpcClient.Close()
	}
	return nil
}

func ConnectBFNToNode(rpcClient *litrpclit.LndcRpcClient, node string) (uint32, error) {
	bfnConnectMutex.Lock()

	conns1, err := commands.ListConnections(rpcClient)
	if err != nil {
		return 0, err
	}

	retries := 0
	for {
		_, err = commands.Connect(rpcClient, fmt.Sprintf("%s@%s:2448", NodeAddresses[node], node))
		if err == nil {
			break
		} else {
			retries++
		}
		if retries > 5 {
			return 0, err
		}
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

func ConnectAndFund(cli *client.Client, nodeName string) error {
	logging.Info.Printf("Connecting new lit node %s to the big fat node...\n", nodeName)

	time.Sleep(time.Second * 2)

	// Connect to the node (this will block until it's available - since it has to sync blocks and stuff)
	lndc, err := GetLndcRpc(cli, nodeName, true)
	if err != nil {
		logging.Error.Printf("Error connecting to new node %s: %s\n", nodeName, err.Error())
		return err
	}

	logging.Info.Printf("Authorizing admin panel key on new node %s\n", nodeName)
	adminPanelKey, err := GetAdminPanelKey()
	if err != nil {
		logging.Error.Printf("Error getting admin panel key: %s\n", err.Error())
		return err
	}
	adminPanelPubKey := [33]byte{}
	copy(adminPanelPubKey[:], adminPanelKey.PubKey().SerializeCompressed())

	err = litrpc.RcAuth(lndc, adminPanelPubKey, true)
	if err != nil {
		logging.Error.Printf("ConnectAndFund RcAuth error: %s", err.Error())
		return err
	}

	lndc.Close()

	rpcClient, err := GetLndcRpc(cli, "litdemobigfatnode", false)
	if err != nil {
		logging.Error.Printf("Error connecting to BFN: %s\n", err.Error())
		return err
	}

	// Connect the BFN to the new node
	peerIdx, err := ConnectBFNToNode(rpcClient, nodeName)
	if err != nil {
		logging.Error.Printf("Error connecting new node to BFN: %s\n", err.Error())
		return err
	}

	for _, cd := range coindaemons.CoinDaemons {
		logging.Info.Printf("Funding %s with %s\n", nodeName, cd.ContainerName)
		reply, err := commands.Fund(rpcClient, peerIdx, cd.LitCoinType, cd.NodeChannelCapacity, cd.NodeChannelInitialSend)
		if err != nil {
			return err
		}
		logging.Info.Printf("Funded %s - %s\n", cd.ContainerName, reply.Status)
	}

	logging.Info.Println("Funding done - mining block")
	err = coindaemons.MineBlock()
	if err != nil {
		return err
	}

	return nil
}
