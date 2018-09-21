package commands

import "github.com/mit-dci/lit/litrpc"

type TxoInfo struct {
	OutPoint string
	Amt      int64
	Height   int32
	Delay    int32
	CoinType string
	Witty    bool

	KeyPath string
}
type TxoListReply struct {
	Txos []TxoInfo
}

func ListUtxos(c *litrpc.LndcRpcClient) (*TxoListReply, error) {
	reply := new(TxoListReply)
	err := c.Call("LitRPC.TxoList", nil, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
