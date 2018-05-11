package docker

import (
	"context"

	"github.com/gertjaap/lit-demo-setup/admin-api/logging"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var NetworkID string

func InitNetwork(cli *client.Client) error {
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}

	found := false
	for _, n := range networks {
		if n.Name == "lit-demo" {
			found = true
			NetworkID = n.ID
			logging.Info.Print("Network [lit-demo] exists")
			break

		}
	}

	if !found {
		logging.Info.Print("Network [lit-demo] does not exist, creating...")
		res, err := cli.NetworkCreate(context.Background(), "lit-demo", types.NetworkCreate{})
		if err != nil {
			return err
		}
		NetworkID = res.ID
		logging.Info.Print("Network [lit-demo] created")
	}

	c, err := GetAdminPanelContainer(cli)
	if err != nil {
		return err
	}

	found = false
	for k := range c.NetworkSettings.Networks {
		if k == "lit-demo" {
			found = true
		}
	}

	if !found {
		logging.Info.Print("Admin container is not in network [lit-demo], connecting...")
		err = cli.NetworkConnect(context.Background(), NetworkID, c.ID, nil)
		if err != nil {
			return err
		}
		logging.Info.Print("Connected admin panel container to [lit-demo] network")
	}

	return nil
}
