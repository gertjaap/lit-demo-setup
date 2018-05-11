package docker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/gertjaap/lit-demo-setup/admin-api/constants"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var NodeAddresses = map[string]string{}
var creationMutex sync.Mutex

func LitNodes(cli *client.Client) ([]types.Container, error) {
	filteredContainers := []types.Container{}

	knownLitImages := constants.KnownImages()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return filteredContainers, err
	}

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
		for _, ki := range knownLitImages {
			if c.ImageID == ki.ImageID {
				filteredContainers = append(filteredContainers, c)
			}
		}

	}

	return filteredContainers, nil
}

func GetAdminPanelContainer(cli *client.Client) (types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return types.Container{}, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/lit-demo-adminpanel" {
				return c, nil
			}
		}
	}

	return types.Container{}, fmt.Errorf("Admin panel container was not found. Are we not running in Docker?")
}

func HostDataDir(cli *client.Client) (string, error) {
	c, err := GetAdminPanelContainer(cli)
	if err != nil {
		return "", err
	}
	for _, m := range c.Mounts {
		if m.Destination == "/data" {
			return m.Source, nil
		}
	}

	return "", fmt.Errorf("Could not find data directory mount")
}

func BootstrapLitData(idx int) error {
	litPath := path.Join("/data", fmt.Sprintf("lit%02d", idx))
	os.MkdirAll(litPath, 0777)
	key32 := new([32]byte)
	rand.Read(key32[:])
	keyHex := []byte(hex.EncodeToString(key32[:]))
	keyfilePath := path.Join(litPath, "privkey.hex")
	confPath := path.Join(litPath, "lit.conf")
	err := ioutil.WriteFile(keyfilePath, keyHex, 0666)
	if err != nil {
		return err
	}
	conf := []byte("reg=litbtcregtest\nrpchost=0.0.0.0")
	err = ioutil.WriteFile(confPath, conf, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetNewNodeIndex(nodes []types.Container) int {
	maxIndex := 0
	for _, n := range nodes {
		if len(n.Names[0]) == 15 {
			idx, err := strconv.ParseInt(n.Names[0][13:], 10, 32)
			if err == nil {
				if idx > int64(maxIndex) {
					maxIndex = int(idx)
				}
			}
		}
	}

	return maxIndex + 1
}

func DropLitNode(cli *client.Client, name string) error {
	nodes, err := LitNodes(cli)
	if err != nil {
		return err
	}

	containerToDrop := types.Container{ID: "undefined"}

	for _, n := range nodes {
		if n.Names[0][1:] == name {
			containerToDrop = n
		}
	}

	if containerToDrop.ID != "undefined" {
		logging.Info.Println("Found container to drop, dropping...")
		cli.ContainerRemove(context.Background(), containerToDrop.ID, types.ContainerRemoveOptions{Force: true})
		return nil
	}

	logging.Error.Println("Container not found, returning error")
	return fmt.Errorf("Invalid container %s", name)
}

func NewLitNode(cli *client.Client) (models.LitNode, error) {
	newNode := models.LitNode{}
	creationMutex.Lock()
	nodes, err := LitNodes(cli)
	if err != nil {
		return newNode, err
	}

	idx := GetNewNodeIndex(nodes)

	containerConfig := new(container.Config)
	containerConfig.Image = constants.DefaultImage().ImageID
	containerConfig.ExposedPorts = nat.PortSet{
		"2448/tcp": struct{}{},
		"8001/tcp": struct{}{},
	}
	hostDataDir, err := HostDataDir(cli)
	if err != nil {
		return newNode, err
	}
	hostConfig := new(container.HostConfig)
	dataDir := path.Join(hostDataDir, fmt.Sprintf("lit%02d", idx))
	hostConfig.Binds = []string{fmt.Sprintf("%s:/root/.lit", dataDir)}

	hostRpcPort := 51000 + idx
	hostLitPort := 52000 + idx

	hostConfig.PortBindings = nat.PortMap{
		"2448/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(hostLitPort),
			},
		},
		"8001/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(hostRpcPort),
			},
		},
	}

	err = BootstrapLitData(idx)
	if err != nil {
		return newNode, err
	}

	newNodeName := fmt.Sprintf("lit-demo-lit%02d", idx)
	cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, newNodeName)
	if err != nil {
		return newNode, err
	}

	err = cli.NetworkConnect(context.Background(), NetworkID, cbody.ID, nil)
	if err != nil {
		return newNode, err
	}

	err = cli.ContainerStart(context.Background(), cbody.ID, types.ContainerStartOptions{})
	if err != nil {
		return newNode, err
	}
	creationMutex.Unlock()

	// Wait for the node to be available

	newNode.Name = newNodeName
	newNode.PublicLitPort = hostLitPort
	newNode.PublicRpcPort = hostRpcPort

	return newNode, nil
}
