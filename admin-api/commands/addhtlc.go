package commands

import "github.com/mit-dci/lit/litrpc"

// ------------------------- HTLCs
type AddHTLCArgs struct {
	ChanIdx  uint32
	Amt      int64
	LockTime uint32
	RHash    [32]byte
	Data     [32]byte
}
type AddHTLCReply struct {
	StateIndex uint64
	HTLCIndex  uint32
}

func AddHTLC(c *litrpc.LndcRpcClient, chanIdx uint32, amount int64, lockTime uint32, RHash, Data [32]byte) (*AddHTLCReply, error) {
	args := new(AddHTLCArgs)
	args.ChanIdx = chanIdx
	args.Amt = amount
	args.LockTime = lockTime
	copy(args.RHash[:], RHash[:])
	copy(args.Data[:], Data[:])

	reply := new(AddHTLCReply)
	err := c.Call("LitRPC.AddHTLC", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
