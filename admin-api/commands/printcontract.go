package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	White = color.New(color.FgHiWhite).SprintFunc()
)

func PrintContract(c *DlcContract) {
	fmt.Fprintf(color.Output, "%-30s : %d\n", White("Index"), c.Idx)
	fmt.Fprintf(color.Output, "%-30s : [%x...%x...%x]\n",
		White("Oracle public key"),
		c.OracleA[:2], c.OracleA[15:16], c.OracleA[31:])
	fmt.Fprintf(color.Output, "%-30s : [%x...%x...%x]\n",
		White("Oracle R-point"), c.OracleR[:2],
		c.OracleR[15:16], c.OracleR[31:])
	fmt.Fprintf(color.Output, "%-30s : %s\n",
		White("Settlement time"),
		time.Unix(int64(c.OracleTimestamp), 0).UTC().Format(time.UnixDate))
	fmt.Fprintf(color.Output, "%-30s : %d\n",
		White("Funded by us"), c.OurFundingAmount)
	fmt.Fprintf(color.Output, "%-30s : %d\n",
		White("Funded by peer"), c.TheirFundingAmount)
	fmt.Fprintf(color.Output, "%-30s : %d\n",
		White("Coin type"), c.CoinType)

	peer := "None"
	if c.PeerIdx > 0 {
		peer = fmt.Sprintf("Peer %d", c.PeerIdx)
	}

	fmt.Fprintf(color.Output, "%-30s : %s\n", White("Peer"), peer)

	status := "Draft"
	switch c.Status {
	case ContractStatusActive:
		status = "Active"
	case ContractStatusClosed:
		status = "Closed"
	case ContractStatusOfferedByMe:
		status = "Sent offer, awaiting reply"
	case ContractStatusOfferedToMe:
		status = "Received offer, awaiting reply"
	case ContractStatusAccepted:
		status = "Accepted"
	case ContractStatusDeclined:
		status = "Declined"
	}

	fmt.Fprintf(color.Output, "%-30s : %s\n\n", White("Status"), status)
}
