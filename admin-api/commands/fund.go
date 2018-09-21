package commands

import (
	"time"

	"github.com/mit-dci/lit/litrpc"
)

type FundArgs struct {
	Peer        uint32 // who to make the channel with
	CoinType    uint32 // what coin to use
	Capacity    int64  // later can be minimum capacity
	Roundup     int64  // ignore for now; can be used to round-up capacity
	InitialSend int64  // Initial send of -1 means "ALL"
	Data        [32]byte
}

func Fund(c *litrpc.LndcRpcClient, peerIdx, coinType uint32, amount, initialSend int64) (*StatusReply, error) {
	args := new(FundArgs)
	args.Peer = peerIdx
	args.CoinType = coinType
	args.Capacity = amount
	args.InitialSend = initialSend

	reply := new(StatusReply)

	ch := make(chan error, 1)
	go func() { ch <- c.Call("LitRPC.FundChannel", args, &reply) }()
	select {
	case err := <-ch:
		return reply, err
	case <-time.After(time.Second * 10):
		return reply, nil
	}
}
