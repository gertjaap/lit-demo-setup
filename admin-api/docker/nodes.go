package docker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemon"
	"github.com/gertjaap/lit-demo-setup/admin-api/constants"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/models"
	"github.com/mit-dci/lit/btcutil/hdkeychain"
	"github.com/mit-dci/lit/coinparam"
	"github.com/mit-dci/lit/litrpc"
	"github.com/mit-dci/lit/lnutil"
	"github.com/mit-dci/lit/portxo"
)

var NodeAddresses = map[string]string{}
var creationMutex sync.Mutex

func SetLitImage(cli *client.Client) error {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return err
	}

	found := false
	for _, i := range images {
		for _, t := range i.RepoTags {
			if t == "lit:latest" {
				constants.SetDefaultImage(i.ID)
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("Could not find image tagged 'lit:latest'")
	}
	return nil
}

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
	err := ioutil.WriteFile(keyfilePath, keyHex, 0666)
	if err != nil {
		return err
	}

	configString := ""

	for _, cd := range coindaemon.CoinDaemons {
		configString += fmt.Sprintf("%s=%s:%d\n", cd.LitConfigPrefix, cd.ContainerName, cd.P2PPort)
	}

	configString += "tracker=http://litdemotracker:46580\n"

	conf := []byte(configString)
	confPath := path.Join(litPath, "lit.conf")
	err = ioutil.WriteFile(confPath, conf, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetNewNodeIndex(nodes []types.Container) int {
	maxIndex := 0
	for _, n := range nodes {
		if len(n.Names[0]) == 13 {
			idx, err := strconv.ParseInt(n.Names[0][11:], 10, 32)
			if err == nil {
				if idx > int64(maxIndex) {
					maxIndex = int(idx)
				}
			}
		}
		if n.Names[0] == "/litdemobigfatnode" {
			if maxIndex == 0 {
				maxIndex = 1
			}
		}
	}

	return maxIndex + 1
}

func GetLitNodeDataDir(cli *client.Client, name string) (string, error) {
	// Find the lit node
	nodes, err := LitNodes(cli)
	if err != nil {
		return "", err
	}

	container := types.Container{ID: "undefined"}

	for _, n := range nodes {
		if n.Names[0][1:] == name {
			container = n
		}
	}

	// Find the directory where the privkey.hex for this container is stored
	hostDataDir, err := HostDataDir(cli)
	if err != nil {
		return "", err
	}

	dataDir := "/data"
	for _, m := range container.Mounts {
		if m.Destination == "/root/.lit" {
			dataDir += strings.Replace(m.Source, hostDataDir, "", 1)
			break
		}
	}

	return dataDir, nil
}

func GetAddress(cli *client.Client, name string) (string, error) {
	dataDir, err := GetLitNodeDataDir(cli, name)
	if err != nil {
		return "", err
	}

	keyFilePath := filepath.Join(dataDir, "privkey.hex")
	privKey, err := lnutil.ReadKeyFile(keyFilePath)
	if err != nil {
		return "", err
	}
	rootPrivKey, err := hdkeychain.NewMaster(privKey[:], &coinparam.TestNet3Params)
	if err != nil {
		return "", err
	}

	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 513 | 1<<31
	kg.Step[2] = 9 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = 0 | 1<<31
	localIDPriv, err := kg.DerivePrivateKey(rootPrivKey)
	if err != nil {
		logging.Error.Printf(err.Error())
	}
	var localIDPub [33]byte
	copy(localIDPub[:], localIDPriv.PubKey().SerializeCompressed())

	return lnutil.LitAdrFromPubkey(localIDPub), nil
}

func GetLndcRpc(cli *client.Client, name string) (*litrpc.LndcRpcClient, error) {
	dataDir, err := GetLitNodeDataDir(cli, name)
	if err != nil {
		return nil, err
	}
	return NewLndcFromHostNameAndDataDir(name, dataDir)
}

func NewLndcFromHostNameAndDataDir(hostName, dataDir string) (*litrpc.LndcRpcClient, error) {
	keyFilePath := filepath.Join(dataDir, "privkey.hex")
	privKey, err := lnutil.ReadKeyFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	rootPrivKey, err := hdkeychain.NewMaster(privKey[:], &coinparam.TestNet3Params)
	if err != nil {
		return nil, err
	}

	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 513 | 1<<31
	kg.Step[2] = 9 | 1<<31
	kg.Step[3] = 1 | 1<<31
	kg.Step[4] = 0 | 1<<31
	key, err := kg.DerivePrivateKey(rootPrivKey)
	if err != nil {
		return nil, err
	}

	kg.Step[3] = 0 | 1<<31
	localIDPriv, err := kg.DerivePrivateKey(rootPrivKey)
	if err != nil {
		logging.Error.Printf(err.Error())
	}
	var localIDPub [33]byte
	copy(localIDPub[:], localIDPriv.PubKey().SerializeCompressed())

	adr := fmt.Sprintf("%s@%s:%d", lnutil.LitAdrFromPubkey(localIDPub), hostName, 2448)
	localIDPriv = nil

	retries := 0
	var ret *litrpc.LndcRpcClient
	for true {
		ret, err = litrpc.NewLndcRpcClient(adr, key)
		if err != nil {
			retries++
			if retries > 10 {
				return nil, err
			}
		} else {
			break
		}
		time.Sleep(time.Second * 1)
	}

	return ret, nil
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
	}
	hostDataDir, err := HostDataDir(cli)
	if err != nil {
		return newNode, err
	}
	hostConfig := new(container.HostConfig)
	dataDir := path.Join(hostDataDir, fmt.Sprintf("lit%02d", idx))
	hostConfig.Binds = []string{fmt.Sprintf("%s:/root/.lit", dataDir)}

	hostLitPort := 52000 + idx

	hostConfig.PortBindings = nat.PortMap{
		"2448/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(hostLitPort),
			},
		},
	}

	err = BootstrapLitData(idx)
	if err != nil {
		return newNode, err
	}

	newNodeName := fmt.Sprintf("litdemolit%02d", idx)
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
	newNode.ID = cbody.ID
	newNode.PublicLitPort = hostLitPort

	return newNode, nil
}
