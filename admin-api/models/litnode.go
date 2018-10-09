package models

import (
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
)

type LitNode struct {
	Name          string
	ID            string
	Balances      map[uint32]int64
	Address       string
	PublicLitPort int
	PublicRpcPort int
	Error         bool
	ErrorDetails  string
	TrackerOK     bool
	TrackerIP     string
	Channels      []commands.ChannelInfo
}
