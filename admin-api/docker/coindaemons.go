package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
)

func WriteConf(coinDaemon coindaemons.CoinDaemon) error {
	confDir := path.Join("/data", coinDaemon.DataSubFolderOnHost)
	os.MkdirAll(confDir, 0777)
	confPath := path.Join(confDir, coinDaemon.ConfigName)
	conf := []byte("regtest=1\nprinttoconsole=1\ndnsseed=0\nupnp=0\nrpcuser=lit\nrpcpassword=lit\nserver=1\ntxindex=1\nrpcallowip=0.0.0.0/0")
	err := ioutil.WriteFile(confPath, conf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InitCoinDaemons(cli *client.Client) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, cd := range coindaemons.CoinDaemons {
		if cd.ImageID == "" {
			cd.ImageID, err = GetImage(cli, cd.ImageName)
			if err != nil {
				return err
			}
		}
		logging.Info.Printf("Checking if coin daemon %s is running\n", cd.ContainerName)
		found := false
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
				if n[1:] == cd.ContainerName && c.ImageID != cd.ImageID {
					cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
				}
			}
			if c.ImageID == cd.ImageID {
				found = true
			}
		}

		if !found {
			logging.Info.Printf("Coin daemon %s is not running, creating...\n", cd.ContainerName)

			err = WriteConf(cd)
			if err != nil {
				return err
			}
			containerConfig := new(container.Config)
			containerConfig.Image = cd.ImageID
			containerConfig.Volumes = map[string]struct{}{
				cd.DataFolderInContainer: {}}
			containerConfig.ExposedPorts = nat.PortSet{
				nat.Port(fmt.Sprintf("%d/tcp", cd.P2PPort)): struct{}{},
				nat.Port(fmt.Sprintf("%d/tcp", cd.RPCPort)): struct{}{},
			}
			if cd.Command != nil {
				containerConfig.Cmd = cd.Command
			}

			hostDataDir, err := HostDataDir(cli)
			if err != nil {
				return err
			}
			hostConfig := new(container.HostConfig)

			hostConfig.PortBindings = nat.PortMap{
				nat.Port(fmt.Sprintf("%d/tcp", cd.P2PPort)): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", cd.P2PPort),
					},
				},
			}

			dataDir := path.Join(hostDataDir, cd.DataSubFolderOnHost)
			hostConfig.Binds = []string{fmt.Sprintf("%s:%s", dataDir, cd.DataFolderInContainer)}
			cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, cd.ContainerName)
			if err != nil {
				return err
			}

			err = cli.NetworkConnect(context.Background(), NetworkID, cbody.ID, nil)
			if err != nil {
				return err
			}

			logging.Info.Printf("Coin daemon %s created, starting...\n", cd.ContainerName)
			err = cli.ContainerStart(context.Background(), cbody.ID, types.ContainerStartOptions{})
			if err != nil {
				return err
			}

			logging.Info.Printf("Coin daemon %s started, waiting for it to boot...\n", cd.ContainerName)

			cd.WaitReady()

			logging.Info.Printf("Coin daemon %s booted, mining first %d blocks...\n", cd.ContainerName, cd.InitialBlocks)
			err = cd.MineBlocks(cd.InitialBlocks)
			if err != nil {
				return err
			}

			logging.Info.Printf("Coin daemon %s operational\n", cd.ContainerName)
		} else {
			logging.Info.Printf("Coin daemon %s is already running\n", cd.ContainerName)
		}
	}
	return nil
}
