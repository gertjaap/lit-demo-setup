package routes

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gorilla/mux"
)

func NodeLogsHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()
	vars := mux.Vars(r)

	containerID, err := docker.GetLitNodeContainerByName(cli, vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log, err := cli.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{Tail: "1000", ShowStdout: true, ShowStderr: true, Timestamps: false})
	w.Header().Set("Content-Type", "text/plain")
	_, err = stdcopy.StdCopy(w, w, log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
