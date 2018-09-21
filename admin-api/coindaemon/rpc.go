package coindaemon

import (
	"fmt"
	"time"

	btcdchaincfg "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	btcdbtcutil "github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/gertjaap/lit-demo-setup/admin-api/logging"
)

func GetRpcClient(coinDaemon CoinDaemon) (*rpcclient.Client, error) {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", coinDaemon.ContainerName, coinDaemon.RPCPort),
		User:         "lit",
		Pass:         "lit",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return client, nil

}

func (coinDaemon CoinDaemon) WaitReady() error {
	logging.Info.Printf("Waiting for %s RPC to be ready...\n", coinDaemon.ContainerName)
	client, err := GetRpcClient(coinDaemon)
	if err != nil {
		return err
	}

	i := 0
	for true {
		blocks, err := client.GetBlockCount()
		if err == nil {
			logging.Info.Printf("%s RPC to be ready (%d blocks)\n", coinDaemon.ContainerName, blocks)
			break
		} else {
			i++
			logging.Info.Printf("%s RPC not ready, waiting (%d sec)...\r", coinDaemon.ContainerName, i)
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

func (coinDaemon CoinDaemon) MineBlocks(num uint32) error {
	client, err := GetRpcClient(coinDaemon)
	if err != nil {
		return err
	}

	_, err = client.Generate(num)
	if err != nil {
		return err
	}

	return nil
}

// decodeSegWitAddress parses a bech32 encoded segwit address string and
// returns the witness version and witness program byte representation.
func decodeSegWitAddress(address string) (byte, []byte, error) {
	// Decode the bech32 encoded address.
	_, data, err := bech32.Decode(address)
	if err != nil {
		return 0, nil, err
	}

	// The first byte of the decoded address is the witness version, it must
	// exist.
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("no witness version")
	}

	// ...and be <= 16.
	version := data[0]
	if version > 16 {
		return 0, nil, fmt.Errorf("invalid witness version: %v", version)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original witness program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return 0, nil, err
	}

	// The regrouped data must be between 2 and 40 bytes.
	if len(regrouped) < 2 || len(regrouped) > 40 {
		return 0, nil, fmt.Errorf("invalid data length")
	}

	// For witness version 0, address MUST be exactly 20 or 32 bytes.
	if version == 0 && len(regrouped) != 20 && len(regrouped) != 32 {
		return 0, nil, fmt.Errorf("invalid data length for witness "+
			"version 0: %v", len(regrouped))
	}

	return version, regrouped, nil
}

func (coinDaemon CoinDaemon) SendCoins(addr string, amt int64) error {
	client, err := GetRpcClient(coinDaemon)
	if err != nil {
		return err
	}

	_, pk, err := decodeSegWitAddress(addr)
	if err != nil {
		return err
	}

	realAddr, err := btcdbtcutil.NewAddressWitnessPubKeyHash(pk, &btcdchaincfg.Params{
		Bech32HRPSegwit:  coinDaemon.CoinParams.Bech32Prefix,
		PubKeyHashAddrID: coinDaemon.CoinParams.PubKeyHashAddrID,
		ScriptHashAddrID: coinDaemon.CoinParams.ScriptHashAddrID,
		PrivateKeyID:     coinDaemon.CoinParams.PrivateKeyID,
		HDPrivateKeyID:   coinDaemon.CoinParams.HDPrivateKeyID,
		HDPublicKeyID:    coinDaemon.CoinParams.HDPublicKeyID,
	})

	if err != nil {

		return err
	}

	_, err = client.SendToAddress(realAddr, btcdbtcutil.Amount(amt))
	if err != nil {
		return err
	}

	return nil
}

func (coinDaemon CoinDaemon) CheckTx(txidHex string) (bool, error) {
	client, err := GetRpcClient(coinDaemon)
	if err != nil {
		return false, err
	}
	txHash, err := chainhash.NewHashFromStr(txidHex)
	if err != nil {
		return false, err
	}
	tx, err := client.GetRawTransaction(txHash)
	if err != nil {
		return false, err
	}
	return (tx.Hash().String() == txidHex), nil
}
