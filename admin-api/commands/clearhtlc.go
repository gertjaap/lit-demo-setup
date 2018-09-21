package commands

import "github.com/mit-dci/lit/litrpc"

type ClearHTLCArgs struct {
	ChanIdx uint32
	HTLCIdx uint32
	R       [16]byte
	Data    [32]byte
}
type ClearHTLCReply struct {
	StateIndex uint64
}

func ClearHTLC(c *litrpc.LndcRpcClient, channelIdx, htlcIndex uint32, R [16]byte, Data [32]byte) (*ClearHTLCReply, error) {
	args := new(ClearHTLCArgs)
	args.ChanIdx = channelIdx
	args.HTLCIdx = htlcIndex
	copy(args.R[:], R[:])
	copy(args.Data[:], Data[:])

	reply := new(ClearHTLCReply)
	err := c.Call("LitRPC.ClearHTLC", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
