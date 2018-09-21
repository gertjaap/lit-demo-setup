package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// DlcOracleRPointResponse is the response format for the REST API that returns
// the R-point
type DlcOracleRPointResponse struct {
	RHex string `json:"R"`
}

func (o *DlcOracle) FetchRPoint(datafeedId, timestamp uint64) ([33]byte, error) {
	var rPoint [33]byte
	if len(o.Url) == 0 {
		return rPoint, fmt.Errorf("Oracle was not imported from the web -" +
			" cannot fetch R point. Enter manually using the" +
			" [dlc contract setrpoint] command")
	}

	url := fmt.Sprintf("%s/api/rpoint/%d/%d", o.Url, datafeedId, timestamp)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return rPoint, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return rPoint, err
	}
	defer resp.Body.Close()

	var response DlcOracleRPointResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return rPoint, err
	}

	R, err := hex.DecodeString(response.RHex)
	if err != nil {
		return rPoint, err
	}

	copy(rPoint[:], R[:])
	return rPoint, nil

}

type DlcOracleSignatureResponse struct {
	Value        int64  `json:"value"`
	SignatureHex string `json:"signature"`
}

func (o *DlcOracle) FetchSignature(rPoint [33]byte) (int64, [32]byte, error) {
	var sig [32]byte
	if len(o.Url) == 0 {
		return 0, sig, fmt.Errorf("Oracle was not imported from the web -" +
			" cannot fetch signature. Enter manually using the" +
			" [dlc contract settle] command")
	}

	url := fmt.Sprintf("%s/api/publication/%s", o.Url, hex.EncodeToString(rPoint[:]))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, sig, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, sig, err
	}
	defer resp.Body.Close()

	var response DlcOracleSignatureResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, sig, err
	}

	s, err := hex.DecodeString(response.SignatureHex)
	if err != nil {
		return 0, sig, err
	}

	copy(sig[:], s[:])
	return response.Value, sig, nil

}
