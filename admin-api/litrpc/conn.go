package litrpc

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"

	"golang.org/x/net/websocket"
)

func ConnectWithRetry(rpcUrl string) (*websocket.Conn, *rpc.Client, error) {
	wsConn, err := websocket.Dial(rpcUrl, "", "http://127.0.0.1/")
	tries := 0
	for {
		if err == nil {
			break
		}
		logging.Info.Printf("Error connecting to %s (%s), retrying %d more times...", rpcUrl, err.Error(), 10-tries)
		if tries > 10 {
			return nil, nil, fmt.Errorf("Could not connect to newly spawn node")
		}
		time.Sleep(time.Millisecond * 500)
		wsConn, err = websocket.Dial(rpcUrl, "", "http://127.0.0.1/")
		tries++
	}

	rpcCon := jsonrpc.NewClient(wsConn)
	return wsConn, rpcCon, nil
}

func Connect(rpcUrl string) (*websocket.Conn, *rpc.Client, error) {
	wsConn, err := websocket.Dial(rpcUrl, "", "http://127.0.0.1/")
	if err != nil {
		return nil, nil, err
	}

	rpcCon := jsonrpc.NewClient(wsConn)
	return wsConn, rpcCon, nil
}
