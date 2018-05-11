package models

type LitNode struct {
	Name          string
	Balances      map[uint32]int64
	Address       string
	PublicLitPort int
	PublicRpcPort int
}
