package main

import (
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemon"
	"github.com/gertjaap/lit-demo-setup/admin-api/docker"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
	"github.com/gertjaap/lit-demo-setup/admin-api/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.HandleFunc("/api/nodes/list", routes.ListNodesHandler)
	r.HandleFunc("/api/nodes/new", routes.NewNodeHandler)
	r.HandleFunc("/api/nodes/delete/{id}", routes.DeleteNodeHandler)
	r.HandleFunc("/api/chain/height", routes.BlockHeightHandler)
	r.HandleFunc("/api/chain/mine", routes.MineBlockHandler)
	r.HandleFunc("/api/redirecttowebui", routes.RedirectToWebUiHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	miner := time.NewTicker(30 * time.Second)
	go func() {
		for range miner.C {
			for _, cd := range coindaemon.CoinDaemons {
				err := cd.MineBlocks(1)
				if err != nil {
					logging.Error.Printf("Could not mine block on %s: %s\n", cd.ContainerName, err)
				}
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
