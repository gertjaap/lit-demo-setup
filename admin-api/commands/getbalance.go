package commands

import "github.com/mit-dci/lit/litrpc"

type CoinBalReply struct {
	CoinType    uint32
	SyncHeight  int32 // height this wallet is synced to
	ChanTotal   int64 // total balance in channels
	TxoTotal    int64 // all utxos
	MatureWitty int64 // confirmed, spendable and witness
	FeeRate     int64 // fee per byte
}

type BalanceReply struct {
	Balances []CoinBalReply
}

func GetBalance(c *litrpc.LndcRpcClient) (*BalanceReply, error) {
	reply := new(BalanceReply)
	err := c.Call("LitRPC.Balance", nil, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
