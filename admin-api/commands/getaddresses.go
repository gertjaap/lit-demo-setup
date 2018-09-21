package commands

import "github.com/mit-dci/lit/litrpc"

type AddressReply struct {
	WitAddresses    []string
	LegacyAddresses []string
}

func GetAddresses(c *litrpc.LndcRpcClient) (*AddressReply, error) {
	reply := new(AddressReply)
	err := c.Call("LitRPC.Address", nil, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
