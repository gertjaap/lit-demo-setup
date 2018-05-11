package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/gertjaap/lit-docker-tester/btc"
)

const BitcoindImageID = "sha256:b26af599c9eceb36145dc033d9782766397e13f6b987f5a0fa91e8e202ef8768"

func WriteBitcoinConf() error {
	os.MkdirAll("/data/bitcoind", 0777)
	confPath := path.Join("/data/bitcoind/bitcoin.conf")
	conf := []byte("regtest=1\ndnsseed=0\nupnp=0\nport=18444\nrpcport=19001\nrpcuser=lit\nrpcpassword=lit\nrpcallowip=0.0.0.0/0\nrpcbind=0.0.0.0\nserver=1\ntxindex=1\n")
	err := ioutil.WriteFile(confPath, conf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InitRegTest(cli *client.Client) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

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
		if c.ImageID == BitcoindImageID {
			found = true
		}
	}

	if !found {
		err = WriteBitcoinConf()
		if err != nil {
			return err
		}
		containerConfig := new(container.Config)
		containerConfig.Image = BitcoindImageID
		containerConfig.Volumes = map[string]struct{}{
			"/bitcoin/.bitcoin": {}}

		hostDataDir, err := HostDataDir(cli)
		if err != nil {
			return err
		}
		hostConfig := new(container.HostConfig)
		dataDir := path.Join(hostDataDir, "bitcoind")
		hostConfig.Binds = []string{fmt.Sprintf("%s:/bitcoin/.bitcoin", dataDir)}
		cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, "litbtcregtest")
		if err != nil {
			return err
		}

		err = cli.NetworkConnect(context.Background(), NetworkID, cbody.ID, nil)
		if err != nil {
			return err
		}

		err = cli.ContainerStart(context.Background(), cbody.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}

		// Wait for Regtest to boot
		time.Sleep(time.Second * 5)

		err = btc.MineBlocks(200)
		if err != nil {
			return err
		}
	}

	return nil
}
