package routes

import (
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/commands"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
)

type ChannelGraphReply struct {
	Graph string
}

func ChannelGraphHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		logging.Error.Printf("ChannelGraphHandler NewEnvClient error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rpcClient, err := docker.GetLndcRpc(cli, "litdemobigfatnode", false)
	if err != nil {
		logging.Error.Printf("ChannelGraphHandler Error connecting to BFN: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var args commands.NoArgs
	reply := new(ChannelGraphReply)
	err = rpcClient.Call("LitRPC.GetChannelMap", args, reply)
	if err != nil {
		logging.Error.Printf("ChannelGraphHandler GetChannelMap Error : %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/png")

	dot := exec.Command("dot", "-Tpng")
	dot.Stdin = strings.NewReader(reply.Graph)
	dot.Stdout = w
	dot.Stderr = os.Stderr
	err = dot.Run()
	if err != nil {
		logging.Error.Printf("ChannelGraphHandler Error running dot: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
