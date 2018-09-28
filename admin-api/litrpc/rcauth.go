package litrpc

import (
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/mit-dci/lit/litrpc"
)

func PendingRCAuthRequests(rpcCon *litrpc.LndcRpcClient) ([][33]byte, error) {
	res, err := commands.PendingRCAuthRequests(rpcCon)
	if err != nil {
		logging.Error.Println(err)
		return [][33]byte{}, err
	}
	return res.PubKeys, nil
}

func RcAuth(rpcCon *litrpc.LndcRpcClient, pubKey [33]byte, auth bool) error {
	_, err := commands.RCAuth(rpcCon, pubKey, auth)
	if err != nil {
		logging.Error.Println(err)
		return err
	}

	return nil
}
