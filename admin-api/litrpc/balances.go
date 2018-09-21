package litrpc

import (
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/mit-dci/lit/litrpc"
)

func GetBalancesFromNode(rpcCon *litrpc.LndcRpcClient) (map[uint32]int64, error) {
	returnVal := map[uint32]int64{}
	res, err := commands.GetBalance(rpcCon)
	if err != nil {

		logging.Error.Println(err)
		return returnVal, err

	}

	for _, b := range res.Balances {
		returnVal[b.CoinType] = b.MatureWitty
	}
	return returnVal, nil
}
