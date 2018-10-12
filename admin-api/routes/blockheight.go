package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gertjaap/lit-demo-setup/admin-api/coindaemons"
)

var blockHeightsLastRefreshed time.Time
var cachedBlockHeights map[string]int64

func BlockHeightHandler(w http.ResponseWriter, r *http.Request) {
	if blockHeightsLastRefreshed.Add(20 * time.Second).Before(time.Now()) { // Cache expired.
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

		cachedBlockHeights = result
		blockHeightsLastRefreshed = time.Now()
	}

	js, err := json.Marshal(cachedBlockHeights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
