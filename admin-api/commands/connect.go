package commands

import "github.com/mit-dci/lit/litrpc"

type ConnectArgs struct {
	LNAddr string
}

func Connect(c *litrpc.LndcRpcClient, addr string) (*StatusReply, error) {
	args := new(ConnectArgs)
	args.LNAddr = addr

	reply := new(StatusReply)
	err := c.Call("LitRPC.Connect", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
