package commands

import "github.com/mit-dci/lit/litrpc"

type ListConnectionsReply struct {
	Connections []PeerInfo
	MyPKH       string
}

type PeerInfo struct {
	PeerNumber uint32
	RemoteHost string
	Nickname   string
}

func ListConnections(c *litrpc.LndcRpcClient) (*ListConnectionsReply, error) {
	args := new(NoArgs)
	reply := new(ListConnectionsReply)
	err := c.Call("LitRPC.ListConnections", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
