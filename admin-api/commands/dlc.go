package commands

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/mit-dci/lit/litrpc"
)

type ListOraclesArgs struct {
	// none
}

type ListOraclesReply struct {
	Oracles []*DlcOracle
}

type DlcOracle struct {
	Idx  uint64   // Index of the oracle for refencing in commands
	A    [33]byte // public key of the oracle
	Name string   // Name of the oracle for display purposes
	Url  string   // Base URL of the oracle, if its REST based (optional)
}

type ImportOracleArgs struct {
	Url  string
	Name string
}

type ImportOracleReply struct {
	Oracle *DlcOracle
}

type AddOracleArgs struct {
	Key  string
	Name string
}

type AddOracleReply struct {
	Oracle *DlcOracle
}

type DlcContractStatus int

const (
	ContractStatusDraft        DlcContractStatus = 0
	ContractStatusOfferedByMe  DlcContractStatus = 1
	ContractStatusOfferedToMe  DlcContractStatus = 2
	ContractStatusDeclined     DlcContractStatus = 3
	ContractStatusAccepted     DlcContractStatus = 4
	ContractStatusAcknowledged DlcContractStatus = 5
	ContractStatusActive       DlcContractStatus = 6
	ContractStatusClosed       DlcContractStatus = 7
)

type DlcContractSettlementSignature struct {
	Outcome   int64    // The oracle value for which transaction these are the signatures
	Signature [64]byte // The signature for the transaction
}

type DlcContract struct {
	Idx                                      uint64                           // Index of the contract for referencing in commands
	TheirIdx                                 uint64                           // Index of the contract on the other peer (so we can reference it in messages)
	PeerIdx                                  uint32                           // Index of the peer we've offered the contract to or received the contract from
	CoinType                                 uint32                           // Coin type
	OracleA, OracleR                         [33]byte                         // Pub keys of the oracle
	OracleTimestamp                          uint64                           // The time we expect the oracle to publish
	Division                                 []DlcContractDivision            // The payout specification
	OurFundingAmount, TheirFundingAmount     int64                            // The amounts either side are funding
	OurChangePKH, TheirChangePKH             [20]byte                         // PKH to which the contracts funding change should go
	OurFundMultisigPub, TheirFundMultisigPub [33]byte                         // Pubkey used in the funding multisig output
	OurPayoutBase, TheirPayoutBase           [33]byte                         // Pubkey to be used in the commit script (combined with oracle pubkey or CSV timeout)
	OurPayoutPKH, TheirPayoutPKH             [20]byte                         // Pubkeyhash to which the contract pays out (directly)
	Status                                   DlcContractStatus                // Status of the contract
	OurFundingInputs, TheirFundingInputs     []DlcContractFundingInput        // Outpoints used to fund the contract
	TheirSettlementSignatures                []DlcContractSettlementSignature // Signatures for the settlement transactions
	FundingOutpoint                          wire.OutPoint                    // The outpoint of the funding TX we want to spend in the settlement - for easier monitoring
}

type DlcContractDivision struct {
	OracleValue int64
	ValueOurs   int64
}

type DlcContractFundingInput struct {
	Outpoint wire.OutPoint
	Value    int64
}

type NewContractArgs struct {
	// empty
}

type NewContractReply struct {
	Contract *DlcContract
}

type ListContractsArgs struct {
	// none
}

type ListContractsReply struct {
	Contracts []*DlcContract
}

type GetContractArgs struct {
	Idx uint64
}

type GetContractReply struct {
	Contract *DlcContract
}

type SetContractOracleArgs struct {
	CIdx uint64
	OIdx uint64
}

type SetContractOracleReply struct {
	Success bool
}

type SetContractDatafeedArgs struct {
	CIdx uint64
	Feed uint64
}

type SetContractDatafeedReply struct {
	Success bool
}

type SetContractRPointArgs struct {
	CIdx   uint64
	RPoint [33]byte
}

type SetContractRPointReply struct {
	Success bool
}

type SetContractSettlementTimeArgs struct {
	CIdx uint64
	Time uint64
}

type SetContractSettlementTimeReply struct {
	Success bool
}

type SetContractFundingArgs struct {
	CIdx        uint64
	OurAmount   int64
	TheirAmount int64
}

type SetContractFundingReply struct {
	Success bool
}

type SetContractDivisionArgs struct {
	CIdx             uint64
	ValueFullyOurs   int64
	ValueFullyTheirs int64
}

type SetContractDivisionReply struct {
	Success bool
}

type SetContractCoinTypeArgs struct {
	CIdx     uint64
	CoinType uint32
}

type SetContractCoinTypeReply struct {
	Success bool
}

type OfferContractArgs struct {
	CIdx    uint64
	PeerIdx uint32
}

type OfferContractReply struct {
	Success bool
}

type DeclineContractArgs struct {
	CIdx uint64
}

type DeclineContractReply struct {
	Success bool
}

type AcceptContractArgs struct {
	CIdx uint64
}

type AcceptContractReply struct {
	Success bool
}

type SettleContractArgs struct {
	CIdx        uint64
	OracleValue int64
	OracleSig   [32]byte
}

type SettleContractReply struct {
	Success bool
}

func ImportOracle(c *litrpc.LndcRpcClient, url, name string) (*ImportOracleReply, error) {
	args := new(ImportOracleArgs)
	args.Url = url
	args.Name = name

	reply := new(ImportOracleReply)
	err := c.Call("LitRPC.ImportOracle", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func NewContract(c *litrpc.LndcRpcClient) (*NewContractReply, error) {
	args := new(NewContractArgs)

	reply := new(NewContractReply)
	err := c.Call("LitRPC.NewContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func ListOracles(c *litrpc.LndcRpcClient) (*ListOraclesReply, error) {
	args := new(ListOraclesArgs)

	reply := new(ListOraclesReply)
	err := c.Call("LitRPC.ListOracles", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func AddOracle(c *litrpc.LndcRpcClient, key, name string) (*AddOracleReply, error) {
	args := new(AddOracleArgs)
	args.Key = key
	args.Name = name
	reply := new(AddOracleReply)
	err := c.Call("LitRPC.AddOracle", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func ListContracts(c *litrpc.LndcRpcClient) (*ListContractsReply, error) {
	args := new(ListContractsArgs)

	reply := new(ListContractsReply)
	err := c.Call("LitRPC.ListContracts", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func GetContract(c *litrpc.LndcRpcClient, idx uint64) (*GetContractReply, error) {
	args := new(GetContractArgs)
	args.Idx = idx
	reply := new(GetContractReply)
	err := c.Call("LitRPC.GetContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractOracle(c *litrpc.LndcRpcClient, cIdx, oIdx uint64) (*SetContractOracleReply, error) {
	args := new(SetContractOracleArgs)
	args.CIdx = cIdx
	args.OIdx = oIdx
	reply := new(SetContractOracleReply)
	err := c.Call("LitRPC.SetContractOracle", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractDatafeed(c *litrpc.LndcRpcClient, cIdx, feed uint64) (*SetContractDatafeedReply, error) {
	args := new(SetContractDatafeedArgs)
	args.CIdx = cIdx
	args.Feed = feed
	reply := new(SetContractDatafeedReply)
	err := c.Call("LitRPC.SetContractDatafeed", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractRPoint(c *litrpc.LndcRpcClient, cIdx uint64, rPoint [33]byte) (*SetContractRPointReply, error) {
	args := new(SetContractRPointArgs)
	args.CIdx = cIdx
	args.RPoint = rPoint
	reply := new(SetContractRPointReply)
	err := c.Call("LitRPC.SetContractRPoint", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractSettlementTime(c *litrpc.LndcRpcClient, cIdx, time uint64) (*SetContractSettlementTimeReply, error) {
	args := new(SetContractSettlementTimeArgs)
	args.CIdx = cIdx
	args.Time = time
	reply := new(SetContractSettlementTimeReply)
	err := c.Call("LitRPC.SetContractSettlementTime", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractFunding(c *litrpc.LndcRpcClient, cIdx uint64, ours, theirs int64) (*SetContractFundingReply, error) {
	args := new(SetContractFundingArgs)
	args.CIdx = cIdx
	args.OurAmount = ours
	args.TheirAmount = theirs
	reply := new(SetContractFundingReply)
	err := c.Call("LitRPC.SetContractFunding", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractDivision(c *litrpc.LndcRpcClient, cIdx uint64, allOurs, allTheirs int64) (*SetContractDivisionReply, error) {
	args := new(SetContractDivisionArgs)
	args.CIdx = cIdx
	args.ValueFullyOurs = allOurs
	args.ValueFullyTheirs = allTheirs
	reply := new(SetContractDivisionReply)
	err := c.Call("LitRPC.SetContractDivision", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SetContractCoinType(c *litrpc.LndcRpcClient, cIdx uint64, coinType uint32) (*SetContractCoinTypeReply, error) {
	args := new(SetContractCoinTypeArgs)
	args.CIdx = cIdx
	args.CoinType = coinType
	reply := new(SetContractCoinTypeReply)
	err := c.Call("LitRPC.SetContractCoinType", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func OfferContract(c *litrpc.LndcRpcClient, cIdx uint64, peerIdx uint32) (*OfferContractReply, error) {
	args := new(OfferContractArgs)
	args.CIdx = cIdx
	args.PeerIdx = peerIdx
	reply := new(OfferContractReply)
	err := c.Call("LitRPC.OfferContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func AcceptContract(c *litrpc.LndcRpcClient, cIdx uint64) (*AcceptContractReply, error) {
	args := new(AcceptContractArgs)
	args.CIdx = cIdx
	reply := new(AcceptContractReply)
	err := c.Call("LitRPC.AcceptContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func DeclineContract(c *litrpc.LndcRpcClient, cIdx uint64) (*DeclineContractReply, error) {
	args := new(DeclineContractArgs)
	args.CIdx = cIdx
	reply := new(DeclineContractReply)
	err := c.Call("LitRPC.DeclineContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func SettleContract(c *litrpc.LndcRpcClient, cIdx uint64, settleValue int64, oracleSig [32]byte) (*SettleContractReply, error) {
	args := new(SettleContractArgs)
	args.CIdx = cIdx
	args.OracleValue = settleValue
	args.OracleSig = oracleSig
	reply := new(SettleContractReply)
	err := c.Call("LitRPC.SettleContract", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type DlcFwdOffer struct {
	// Convenience definition for serialization from RPC
	OType uint8
	// Index of the offer
	OIdx uint64
	// Index of the offer on the other peer
	TheirOIdx uint64
	// Index of the peer offering to / from
	PeerIdx uint32
	// Coin type
	CoinType uint32
	// Pub keys of the oracle and the R point used in the contract
	OracleA, OracleR [33]byte
	// time of expected settlement
	SettlementTime uint64
	// amount of funding, in sats, each party contributes
	FundAmt int64
	// slice of my payouts for given oracle prices
	Payouts []DlcContractDivision

	// if true, I'm the 'buyer' of the foward asset (and I'm short bitcoin)
	ImBuyer bool

	// amount of asset to be delivered at settlement time
	// note that initial price is FundAmt / AssetQuantity
	AssetQuantity int64

	// Stores if the offer was accepted. When receiving a matching
	// Contract draft, we can accept that too.
	Accepted bool
}

type NewForwardOfferArgs struct {
	Offer *DlcFwdOffer
}

type NewForwardOfferReply struct {
	Offer *DlcFwdOffer
}

func NewForwardOffer(c *litrpc.LndcRpcClient, offer *DlcFwdOffer) (*NewForwardOfferReply, error) {

	args := new(NewForwardOfferArgs)
	args.Offer = offer
	reply := new(NewForwardOfferReply)
	err := c.Call("LitRPC.NewForwardOffer", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type ListOffersArgs struct {
	// none
}

type ListOffersReply struct {
	Offers []*DlcFwdOffer
}

func ListOffers(c *litrpc.LndcRpcClient) (*ListOffersReply, error) {

	args := new(ListOffersArgs)
	reply := new(ListOffersReply)
	err := c.Call("LitRPC.ListOffers", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type DeclineOfferArgs struct {
	OIdx uint64
}

type DeclineOfferReply struct {
	Success bool
}

func DeclineOffer(c *litrpc.LndcRpcClient, oIdx uint64) (*DeclineOfferReply, error) {

	args := new(DeclineOfferArgs)
	args.OIdx = oIdx
	reply := new(DeclineOfferReply)
	err := c.Call("LitRPC.DeclineOffer", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type AcceptOfferArgs struct {
	OIdx uint64
}

type AcceptOfferReply struct {
	Success bool
}

func AcceptOffer(c *litrpc.LndcRpcClient, oIdx uint64) (*AcceptOfferReply, error) {

	args := new(AcceptOfferArgs)
	args.OIdx = oIdx
	reply := new(AcceptOfferReply)
	err := c.Call("LitRPC.AcceptOffer", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
