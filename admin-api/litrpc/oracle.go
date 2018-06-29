package litrpc

import (
	"net/rpc"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-docker-tester/commands"
)

func ImportOracle(rpcCon *rpc.Client) error {
	_, err := commands.ImportOracle(rpcCon, "http://litoracle:3000/", "LIT Oracle")
	if err != nil {

		logging.Error.Println(err)
		return err

	}

	return nil
}
