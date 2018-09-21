package commands

import "github.com/mit-dci/lit/litrpc"

type ChannelInfo struct {
	OutPoint      string
	CoinType      uint32
	Closed        bool
	Capacity      int64
	MyBalance     int64
	Height        int32  // block height of channel fund confirmation
	StateNum      uint64 // Most recent commit number
	PeerIdx, CIdx uint32
	PeerID        string
	Data          [32]byte
	Pkh           [20]byte
}
type ChannelListReply struct {
	Channels []ChannelInfo
}

func ListChannels(c *litrpc.LndcRpcClient) (*ChannelListReply, error) {
	args := new(ChanArgs)
	reply := new(ChannelListReply)
	err := c.Call("LitRPC.ChannelList", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
