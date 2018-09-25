package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
)

func BlockHeightHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]int64{}

	for _, cd := range coindaemons.CoinDaemons {
		cli, err := coindaemons.GetRpcClient(cd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		blocks, err := cli.GetBlockCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result[cd.DataSubFolderOnHost] = blocks
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
