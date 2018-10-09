package models

import (
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
)

type LitNode struct {
	Name          string
	ID            string
	Balances      []commands.CoinBalReply
	Address       string
	PublicLitPort int
	PublicRpcPort int
	Error         bool
	ErrorDetails  string
	TrackerOK     bool
	TrackerIP     string
	Channels      []commands.ChannelInfo
}
