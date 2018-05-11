package litrpc

import (
	"net/rpc"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-docker-tester/commands"
)

func GetAddressFromNode(rpcCon *rpc.Client) (string, error) {
	res, err := commands.GetListeningPorts(rpcCon)
	if err != nil {
		res, err = commands.Listen(rpcCon, ":2448")
		if err != nil {
			logging.Error.Println(err)
			return "", err
		}
	}

	return res.Adr, nil
}

func GetBlockchainAddress(rpcCon *rpc.Client) (string, error) {
	res, err := commands.GetAddresses(rpcCon)
	if err != nil {
		logging.Error.Println(err)
		return "", err
	}

	return res.WitAddresses[0], nil
}
