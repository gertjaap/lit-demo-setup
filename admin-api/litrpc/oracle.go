package litrpc

import (
	"net/rpc"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-docker-tester/commands"
)

func ImportOracle(rpcCon *rpc.Client) error {
	_, err := commands.ImportOracle(rpcCon, "https://oracle.gertjaap.org/", "Demo Oracle")
	if err != nil {

		logging.Error.Println(err)
		return err

	}

	return nil
}
