package main

import (
	"github.com/balletcrypto/bitcoin-inscription-parser/parser"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Create an RPC client that connects to a bitcoin node
	config := &rpcclient.ConnConfig{
		Host:         "your rpc host",
		User:         "your rpc user",
		Pass:         "your rpc password",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(config, nil)
	if err != nil {
		log.Fatalf("Create rpc client connection to bitcoind node failed, error: %v", err)
	}
	defer client.Shutdown()

	// Get the raw transaction data of the specified tx hash
	txHash := "fe76628c921e7894e4f34f036cd081fc4b21009639d6f4fc12577f59818b35b8"
	hashFromStr, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Fatalf("Get tx hash from string failed, error: %v", err)
	}

	rawTx, err := client.GetRawTransaction(hashFromStr)
	if err != nil {
		log.Fatalf("Get raw tx failed, error: %v", err)
	}
	transactionInscriptions := parser.ParseInscriptionsFromTransaction(rawTx.MsgTx())
	if len(transactionInscriptions) == 0 {
		log.Infof("NO INSCRIPTONS!!!!!")
	}
	for _, v := range transactionInscriptions {
		ins := v
		log.Infof("INCRIPTION txin index: %d, tx in offset: %d, content type: %s, content length: %d",
			ins.TxInIndex, ins.TxInOffset, string(ins.Inscription.ContentType), ins.Inscription.ContentLength)
	}
}
