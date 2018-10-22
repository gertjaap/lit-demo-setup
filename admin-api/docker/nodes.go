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
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
	"github.com/gertjaap/lit-demo-setup/admin-api/constants"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/models"
	"github.com/mit-dci/lit/btcutil/hdkeychain"
	"github.com/mit-dci/lit/coinparam"
	"github.com/mit-dci/lit/crypto/koblitz"
	"github.com/mit-dci/lit/litrpc"
	"github.com/mit-dci/lit/lnutil"
	"github.com/mit-dci/lit/portxo"
)

var NodeAddresses = map[string]string{}
var nodeLndcs = map[string]*litrpc.LndcRpcClient{}
var creationMutex sync.Mutex
var lndcMapMutex sync.Mutex

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
			if n == "/litdemoadminpanel" {
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

	for _, cd := range coindaemons.CoinDaemons {
		configString += fmt.Sprintf("%s=%s:%d\n", cd.LitConfigPrefix, cd.ContainerName, cd.P2PPort)
	}

	configString += "tracker=http://litdemotracker:46580\n"

	conf := []byte(configString)
	confPath := path.Join(litPath, "lit.conf")
	err = ioutil.WriteFile(confPath, conf, 0666)
	if err != nil {
		return err
	}

	// Write the rates.json file. We use 1 BTC = 112 LTC, 118 LTC = 1 BTC, 1 BTC = 6295 USD, 6495 USD = 1 BTC, 1 LTC = 52 USD, 58 USD = 1 LTC
	rates := []byte("{\"257\": [{\"CoinType\": 257,	\"Rate\": 1, \"Reciprocal\": false},{\"CoinType\": 258,	\"Rate\": 118, \"Reciprocal\": true},{\"CoinType\": 262,	\"Rate\": 6495, \"Reciprocal\": true}],\"258\": [{\"CoinType\": 258,	\"Rate\": 1, \"Reciprocal\": false},{\"CoinType\": 257,	\"Rate\": 112, \"Reciprocal\": false},{\"CoinType\": 262,	\"Rate\": 58, \"Reciprocal\": true}],\"262\": [{\"CoinType\": 262,	\"Rate\": 1, \"Reciprocal\": false},{\"CoinType\": 257,	\"Rate\": 6295, \"Reciprocal\": false},{\"CoinType\": 258,	\"Rate\": 52, \"Reciprocal\": false}]}")
	ratesPath := path.Join(litPath, "rates.json")
	err = ioutil.WriteFile(ratesPath, rates, 0666)
	if err != nil {
		return err
	}

	return nil
}

func GetNewNodeIndex(nodes []types.Container) int {
	maxIndex := 0
	for _, n := range nodes {
		if n.Names[0][1:11] == "litdemolit" {
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

var adminPanelKey *koblitz.PrivateKey

func GetAdminPanelKey() (*koblitz.PrivateKey, error) {

	if adminPanelKey == nil {
		key32 := [32]byte{}
		if _, err := os.Stat("/data/adminpanel.key"); os.IsNotExist(err) {
			_, err = rand.Read(key32[:])
			if err != nil {
				adminPanelKey = nil
				return nil, err
			}
			err = ioutil.WriteFile("/data/adminpanel.key", key32[:], 0600)
			if err != nil {
				adminPanelKey = nil
				return nil, err
			}
		} else {
			fileKey, err := ioutil.ReadFile("/data/adminpanel.key")
			if err != nil {
				adminPanelKey = nil
				return nil, err
			}
			copy(key32[:], fileKey)
		}
		adminPanelKey, _ = koblitz.PrivKeyFromBytes(koblitz.S256(), key32[:])
	}
	return adminPanelKey, nil
}

func DropLndcRpc(name string) {
	lndcMapMutex.Lock()
	_, ok := nodeLndcs[name]
	if ok {
		nodeLndcs[name].Close()
		delete(nodeLndcs, name)
	}
	lndcMapMutex.Unlock()
}

func GetLndcRpc(cli *client.Client, name string, useLitAfKey bool) (*litrpc.LndcRpcClient, error) {
	lndcMapMutex.Lock()
	lndc, ok := nodeLndcs[name]
	if ok && !useLitAfKey {
		lndcMapMutex.Unlock()
		return lndc, nil
	}
	dataDir, err := GetLitNodeDataDir(cli, name)
	if err != nil {
		lndcMapMutex.Unlock()
		logging.Error.Printf("Error fetching datadir for %s: %s\n", name, err.Error())
		return nil, err
	}
	rootKey, err := GetRootKeyFromDataDir(dataDir)
	if err != nil {
		lndcMapMutex.Unlock()
		logging.Error.Printf("Error fetching rootKey for %s: %s\n", name, err.Error())
		return nil, err
	}
	localIDPub, err := DeriveNodePub(rootKey)
	adr := fmt.Sprintf("%s@%s:%d", lnutil.LitAdrFromPubkey(localIDPub), name, 2448)

	key, err := GetAdminPanelKey()
	if err != nil {
		lndcMapMutex.Unlock()
		logging.Error.Printf("Error getting adminpanel key: %s\n", err.Error())
		return nil, err
	}
	if useLitAfKey {
		key, err = DeriveLitAfKey(rootKey)
		if err != nil {
			lndcMapMutex.Unlock()
			logging.Error.Printf("Error deriving lit-af key: %s\n", err.Error())
			return nil, err
		}
	}

	retries := 0
	var ret *litrpc.LndcRpcClient
	for true {
		ret, err = litrpc.NewLndcRpcClient(adr, key)
		if err != nil {
			logging.Info.Printf("Error connecting to %s: %s, retrying %d more times\n", adr, err.Error(), 10-retries)
			retries++
			if retries > 10 {
				lndcMapMutex.Unlock()
				return nil, err
			}
		} else {
			break
		}
		time.Sleep(time.Second * 1)
	}

	// Don't cache connections using lit-af key. These should also be closed by the caller.
	if !useLitAfKey {
		nodeLndcs[name] = ret
	}

	lndcMapMutex.Unlock()
	return ret, nil
}

func GetRootKeyFromDataDir(dataDir string) (*hdkeychain.ExtendedKey, error) {
	keyFilePath := filepath.Join(dataDir, "privkey.hex")
	privKey, err := lnutil.ReadKeyFile(keyFilePath)
	if err != nil {
		logging.Error.Printf("NewLndcFromHostNameAndDataDir error in ReadKeyFile %s\n", err.Error())
		return nil, err
	}
	rootPrivKey, err := hdkeychain.NewMaster(privKey[:], &coinparam.TestNet3Params)
	if err != nil {
		logging.Error.Printf("NewLndcFromHostNameAndDataDir error in hdkeychain.NewMaster %s\n", err.Error())
		return nil, err
	}
	return rootPrivKey, nil
}

func DeriveLitAfKey(rootPrivKey *hdkeychain.ExtendedKey) (*koblitz.PrivateKey, error) {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 513 | 1<<31
	kg.Step[2] = 9 | 1<<31
	kg.Step[3] = 1 | 1<<31
	kg.Step[4] = 0 | 1<<31

	key, err := kg.DerivePrivateKey(rootPrivKey)
	if err != nil {
		logging.Error.Printf("NewLndcFromHostNameAndDataDir error in DerivePrivateKey %s\n", err.Error())
		return nil, err
	}
	return key, nil
}

func DeriveNodePub(rootPrivKey *hdkeychain.ExtendedKey) ([33]byte, error) {
	localIDPriv, err := DeriveNodePriv(rootPrivKey)
	if err != nil {
		logging.Error.Printf("NewLndcFromHostNameAndDataDir error in DerivePrivateKey %s\n", err.Error())
		return [33]byte{}, err
	}

	var localIDPub [33]byte
	copy(localIDPub[:], localIDPriv.PubKey().SerializeCompressed())

	return localIDPub, nil
}

func DeriveNodePriv(rootPrivKey *hdkeychain.ExtendedKey) (*koblitz.PrivateKey, error) {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 513 | 1<<31
	kg.Step[2] = 9 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = 0 | 1<<31

	localIDPriv, err := kg.DerivePrivateKey(rootPrivKey)
	if err != nil {
		logging.Error.Printf("NewLndcFromHostNameAndDataDir error in DerivePrivateKey %s\n", err.Error())
		return nil, err
	}
	return localIDPriv, nil
}

func DropLitNode(cli *client.Client, name string) error {
	containerToDrop, err := GetLitNodeContainerByName(cli, name)
	if err != nil {
		return err
	}
	logging.Info.Println("Found container to drop, dropping...")
	cli.ContainerRemove(context.Background(), containerToDrop, types.ContainerRemoveOptions{Force: true})
	return nil
}

func GetLitNodeContainerByName(cli *client.Client, name string) (string, error) {
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

	if container.ID != "undefined" {
		return container.ID, nil
	}

	return "", fmt.Errorf("Container not found")
}

func RestartLitNode(cli *client.Client, name string) error {
	containerToRestart, err := GetLitNodeContainerByName(cli, name)
	if err != nil {
		return err
	}

	logging.Info.Println("Found container to restart, restarting...")
	cli.ContainerRestart(context.Background(), containerToRestart, nil)
	return nil
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

func UpgradeNodes(cli *client.Client) error {
	// Check for containers mounted to the litxx datadirectory that
	// have an old LIT image. We'll restart them with the right image.
	containersToUpgrade := []types.Container{}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	knownLitImages := constants.KnownImages()

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
		imageCorrect := false
		for _, ki := range knownLitImages {
			if c.ImageID == ki.ImageID {
				imageCorrect = true
			}
		}
		if imageCorrect {
			continue
		}

		for _, m := range c.Mounts {
			if m.Destination == "/root/.lit" {
				containersToUpgrade = append(containersToUpgrade, c)
			}
		}

	}

	logging.Info.Printf("Updating %d containers to lit:latest\n", len(containersToUpgrade))
	for _, c := range containersToUpgrade {
		logging.Info.Printf("Dropping %s\n", c.Names[0][1:])
		err := cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}

		logging.Info.Printf("Recreating %s\n", c.Names[0][1:])

		containerConfig := new(container.Config)
		containerConfig.Image = constants.DefaultImage().ImageID
		containerConfig.ExposedPorts = nat.PortSet{
			"2448/tcp": struct{}{},
		}

		hostConfig := new(container.HostConfig)
		for _, m := range c.Mounts {
			if m.Destination == "/root/.lit" {
				hostConfig.Binds = []string{fmt.Sprintf("%s:/root/.lit", m.Source)}
			}
		}

		hostLitPort := 52000
		for _, p := range c.Ports {
			if p.PrivatePort == 2448 {
				hostLitPort = int(p.PublicPort)
			}
		}

		hostConfig.PortBindings = nat.PortMap{
			"2448/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(hostLitPort),
				},
			},
		}

		cbody, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, c.Names[0][1:])
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

		logging.Info.Printf("Completed recreation of %s\n", c.Names[0][1:])

	}

	return nil
}
