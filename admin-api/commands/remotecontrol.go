package commands

import "github.com/mit-dci/lit/litrpc"

type RemoteControlAuthorization struct {
	Allowed bool
}

type RCAuthArgs struct {
	PubKey        [33]byte
	Authorization *RemoteControlAuthorization
}

type RCSendArgs struct {
	PeerIdx uint32
	Msg     []byte
}

type RCPendingAuthRequestsReply struct {
	PubKeys [][33]byte
}

func RCAuth(c *litrpc.LndcRpcClient, pubKey [33]byte, auth bool) (*StatusReply, error) {
	args := new(RCAuthArgs)
	args.Authorization = new(RemoteControlAuthorization)
	args.Authorization.Allowed = auth
	args.PubKey = pubKey

	reply := new(StatusReply)
	err := c.Call("LitRPC.RemoteControlAuth", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func RCSend(c *litrpc.LndcRpcClient, peerIdx uint32, msg []byte) (*StatusReply, error) {
	args := new(RCSendArgs)
	args.PeerIdx = peerIdx
	args.Msg = msg

	reply := new(StatusReply)
	err := c.Call("LitRPC.RemoteControlSend", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func PendingRCAuthRequests(c *litrpc.LndcRpcClient) (*RCPendingAuthRequestsReply, error) {
	args := new(NoArgs)
	reply := new(RCPendingAuthRequestsReply)
	err := c.Call("LitRPC.ListPendingRemoteControlAuthRequests", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
