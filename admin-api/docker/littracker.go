package docker

import (
	"context"
	"fmt"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
)

func InitLitTracker(cli *client.Client) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	mongoImageId, err := GetImage(cli, "mongo")
	if err != nil {
		return err
	}

	littrackerImageId, err := GetImage(cli, "littracker")
	if err != nil {
		return err
	}

	logging.Info.Printf("Checking if lit tracker is running\n")
	trackerFound := false
	mongoFound := false
	trackerId := ""
	for _, c := range containers {
		correctNetwork := false
		for _, n := range c.NetworkSettings.Networks {
			if n.NetworkID == NetworkID {
				correctNetwork = true
			}
		}
		if !correctNetwork {
			continue
		}
		for _, n := range c.Names {
			if n[1:] == "litdemotrackermongo" && c.ImageID != mongoImageId {
				cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
			}
			if n[1:] == "litdemotracker" && c.ImageID != littrackerImageId {
				cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
			}
		}
		if c.ImageID == littrackerImageId {
			trackerFound = true
			trackerId = c.ID
		}
		if c.ImageID == mongoImageId {
			mongoFound = true
		}
	}

	if !mongoFound {
		logging.Info.Printf("Mongo not found, creating...\n")
		containerConfig := new(container.Config)
		containerConfig.Image = mongoImageId
		containerConfig.Volumes = map[string]struct{}{
			"/data/db": {}}

		hostDataDir, err := HostDataDir(cli)
		if err != nil {
			return err
		}
		hostConfig := new(container.HostConfig)
		dataDir := path.Join(hostDataDir, "trackermongo")
		hostConfig.Binds = []string{fmt.Sprintf("%s:%s", dataDir, "/data/db")}
		cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, "litdemotrackermongo")
		if err != nil {
			return err
		}

		err = cli.NetworkConnect(context.Background(), NetworkID, cbody.ID, nil)
		if err != nil {
			return err
		}

		logging.Info.Printf("Mongo starting...\n")
		err = cli.ContainerStart(context.Background(), cbody.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
		logging.Info.Printf("Mongo started\n")
	}

	if !trackerFound {
		logging.Info.Printf("Tracker not found, creating...\n")
		containerConfig := new(container.Config)
		containerConfig.Image = littrackerImageId
		containerConfig.Env = []string{"DB_HOST=litdemotrackermongo"}
		containerConfig.ExposedPorts = nat.PortSet{
			"46580/tcp": struct{}{},
		}
		hostConfig := new(container.HostConfig)

		hostConfig.PortBindings = nat.PortMap{
			"46580/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "46581",
				},
			},
		}

		cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, "litdemotracker")
		if err != nil {
			return err
		}

		err = cli.NetworkConnect(context.Background(), NetworkID, cbody.ID, nil)
		if err != nil {
			return err
		}

		logging.Info.Printf("Tracker starting...\n")
		err = cli.ContainerStart(context.Background(), cbody.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
		logging.Info.Printf("Tracker started\n")
	} else {
		if !mongoFound {
			// If mongo was not found, but tracker was found - restart the tracker
			err = cli.ContainerRestart(context.Background(), trackerId, nil)
			if err != nil {
				return err
			}
			logging.Info.Printf("Tracker restarted due to mongo (re)creation\n")
		}
	}

	return nil
}
