package commands

import "github.com/mit-dci/lit/litrpc"

type ListenArgs struct {
	Port string
}

type ListeningPortsReply struct {
	LisIpPorts []string
	Adr        string
}

func Listen(c *litrpc.LndcRpcClient, port string) (*ListeningPortsReply, error) {
	args := new(ListenArgs)
	args.Port = port

	reply := new(ListeningPortsReply)
	err := c.Call("LitRPC.Listen", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func GetListeningPorts(c *litrpc.LndcRpcClient) (*ListeningPortsReply, error) {
	args := new(NoArgs)

	reply := new(ListeningPortsReply)
	err := c.Call("LitRPC.GetListeningPorts", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
