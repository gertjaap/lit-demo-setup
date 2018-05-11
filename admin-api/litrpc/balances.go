package litrpc

import (
	"net/rpc"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-docker-tester/commands"
)

func GetBalancesFromNode(rpcCon *rpc.Client) (map[uint32]int64, error) {
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
