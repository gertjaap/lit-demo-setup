package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
	"github.com/gertjaap/lit-demo-setup/admin-api/constants"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mit-dci/lit/lnutil"
)

func main() {
	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	logging.Info.Println("Waiting for 5 seconds to get Docker to register us as container...")

	time.Sleep(time.Second * 5)

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	err = constants.InitImages(cli)
	if err != nil {
		panic(err)
	}

	err = docker.InitNetwork(cli)
	if err != nil {
		panic(err)
	}
	err = docker.InitCoinDaemons(cli)
	if err != nil {
		panic(err)
	}
	err = docker.InitLitTracker(cli)
	if err != nil {
		panic(err)
	}
	err = docker.InitBigFatNode(cli)
	if err != nil {
		panic(err)
	}

	err = docker.UpgradeNodes(cli)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/nodes/graph", routes.ChannelGraphHandler)
	r.HandleFunc("/api/nodes/list", routes.ListNodesHandler)
	r.HandleFunc("/api/nodes/new", routes.NewNodeHandler)
	r.HandleFunc("/api/nodes/delete/{id}", routes.DeleteNodeHandler)
	r.HandleFunc("/api/nodes/logs/{id}", routes.NodeLogsHandler)
	r.HandleFunc("/api/nodes/restart/{id}", routes.RestartNodeHandler)
	r.HandleFunc("/api/nodes/pendingauth/{id}", routes.PendingAuthForNodeHandler)
	r.HandleFunc("/api/nodes/auth/{id}/{key}/{yesno}", routes.AuthNodeHandler)
	r.HandleFunc("/api/chain/height", routes.BlockHeightHandler)
	r.HandleFunc("/api/chain/mine", routes.MineBlockHandler)
	r.HandleFunc("/api/redirecttowebui", routes.RedirectToWebUiHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	miner := time.NewTicker(5 * time.Minute)
	go func() {
		for range miner.C {
			for _, cd := range coindaemons.CoinDaemons {
				err := cd.MineBlocks(1)
				if err != nil {
					logging.Error.Printf("Could not mine block on %s: %s\n", cd.ContainerName, err)
				}
			}

		}
	}()

	httpClient := &http.Client{
		Timeout: time.Second * 4, // 4+4 to accomodate the 10s RPC timeout
	}
	resp, err := httpClient.Get("https://api.ipify.org/")
	if err != nil {
		fmt.Printf("Get IPv4 error: %v", err)
		panic(err)
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	myPublicIP := strings.TrimSpace(buf.String())

	announcer := time.NewTicker(60 * time.Second)
	go func() {

		for range announcer.C {
			cli, err := client.NewEnvClient()
			if err != nil {
				logging.Error.Printf("Error getting new env client %s\n", err.Error())
				continue
			}
			defer cli.Close()
			containers, err := docker.LitNodes(cli)
			if err != nil {
				logging.Error.Printf("Error listing nodes: %s\n", err.Error())
				continue
			}

			logging.Info.Printf("Found %d containers to register in the tracker\n", len(containers))

			for _, c := range containers {
				nodeName := c.Names[0][1:]
				logging.Info.Printf("Registering %s in the tracker...\n", nodeName)
				dataDir, err := docker.GetLitNodeDataDir(cli, nodeName)
				if err != nil {
					logging.Error.Printf("Error fetching datadir for %s: %s\n", nodeName, err.Error())
					continue
				}
				rootKey, err := docker.GetRootKeyFromDataDir(dataDir)
				if err != nil {
					logging.Error.Printf("Error fetching rootKey for %s: %s\n", nodeName, err.Error())
					continue
				}
				localIDPub, err := docker.DeriveNodePub(rootKey)
				if err != nil {
					logging.Error.Printf("Error fetching rootKey for %s: %s\n", nodeName, err.Error())
					continue
				}
				localIDPriv, err := docker.DeriveNodePriv(rootKey)
				adr := lnutil.LitAdrFromPubkey(localIDPub)
				publicPort := 0
				for _, p := range c.Ports {
					if p.PrivatePort == 2448 {
						publicPort = int(p.PublicPort)
					}
				}

				liturlIPv4 := fmt.Sprintf("%s:%d", myPublicIP, publicPort)
				liturlIPv6 := "::"
				urlBytes := []byte(liturlIPv4 + liturlIPv6)
				urlHash := sha256.Sum256(urlBytes)
				urlSig, err := localIDPriv.Sign(urlHash[:])
				if err != nil {
					logging.Error.Printf("Error signing announcement for %s: %s\n", nodeName, err.Error())
					continue
				}

				resp, err := http.PostForm("http://litdemotracker:46580/announce",
					url.Values{"ipv4": {liturlIPv4},
						"ipv6": {liturlIPv6},
						"addr": {adr},
						"sig":  {hex.EncodeToString(urlSig.Serialize())},
						"pbk":  {hex.EncodeToString(localIDPub[:])}})

				if err != nil {
					logging.Error.Printf("Error announcing %s: %s\n", nodeName, err.Error())
					continue
				}

				defer resp.Body.Close()
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				logging.Info.Printf("Succesfully announced %s @ %s, result: %s", nodeName, liturlIPv4, buf.String())
			}
		}
	}()

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8000"
	}

	logging.Info.Println("Listening on port %s", portString)

	logging.Error.Fatal(http.ListenAndServe(":"+portString, handlers.CORS(originsOk, headersOk, methodsOk)(logging.WebLoggingMiddleware(r))))
}

type announcement struct {
	ipv4 string
	ipv6 string
	addr string
	sig  string
	pbk  string
}
