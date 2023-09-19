# bitcoin-inscription-parser

bitcoin-inscription-parser is a tool which helps to parse bitcoin inscriptions 
from transactions. Any inscription content wrapped in `OP_FALSE OP_IF … OP_ENDIF`
using data pushes can be correctly parsed.

The tool supports single or multiple inscriptions in all input of the transaction.

# Installation
```
go get github.com/balletcrypto/bitcoin-inscription-parser
```
# Example
```go
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
			ins.TxInIndex, ins.TxInOffset, ins.Inscription.ContentType, ins.Inscription.ContentLength)
	}
}
```
Also shown in examples folder

# Unit tests
```
go test -v script_parser_test.go 
=== RUN   TestScriptWithInscription
=== PAUSE TestScriptWithInscription
=== CONT  TestScriptWithInscription
    script_parser_test.go:75: Find inscription with content type: text/plain;charset=utf-8, content length: 28
    script_parser_test.go:81: test script with inscription: test passed
    script_parser_test.go:75: Find inscription with content type: text/plain;charset=utf-8, content length: 28
    script_parser_test.go:81: test script with unknown tag: test passed
    script_parser_test.go:75: Find inscription with content type: text/plain;charset=utf-8, content length: 28
    script_parser_test.go:81: test script with additions after OP_ENDIF: test passed
--- PASS: TestScriptWithInscription (0.00s)
PASS

=== RUN   TestScriptWithMultipleInscriptions
    script_parser_test.go:121: Find inscription with content type: text/plain;charset=utf-8, content length: 38
    script_parser_test.go:121: Find inscription with content type: text/plain;charset=utf-8, content length: 38
    script_parser_test.go:127: test script with multiple inscriptions: test passed
--- PASS: TestScriptWithMultipleInscriptions (0.00s)
PASS

=== RUN   TestScriptWithIncompleteEnvelopHeader
=== PAUSE TestScriptWithIncompleteEnvelopHeader
=== CONT  TestScriptWithIncompleteEnvelopHeader
    script_parser_test.go:200: test script with no OP_FALSE: test passed
    script_parser_test.go:200: test script with no OP_IF: test passed
    script_parser_test.go:200: test script with no ord: test passed
--- PASS: TestScriptWithIncompleteEnvelopHeader (0.00s)
PASS

=== RUN   TestScriptWithDuplicatedTag
=== PAUSE TestScriptWithDuplicatedTag
=== CONT  TestScriptWithDuplicatedTag
    script_parser_test.go:266: test script with duplicated content type tag: test passed
    script_parser_test.go:266: test script with unknown duplicated tag: test passed
--- PASS: TestScriptWithDuplicatedTag (0.00s)
PASS

=== RUN   TestScriptWithOtherOpcodeBeforeEndIf
    script_parser_test.go:304: test script with other opcode before OP_ENDIF: test passed
--- PASS: TestScriptWithOtherOpcodeBeforeEndIf (0.00s)
PASS

=== RUN   TestScriptWithUnrecognizedEvenTag
    script_parser_test.go:342: test script with unrecognized even tag: test passed
--- PASS: TestScriptWithUnrecognizedEvenTag (0.00s)
PASS

=== RUN   TestScriptWithNoContentType
    script_parser_test.go:369: Find inscription with content type: , content length: 32
    script_parser_test.go:375: test script with no content type: test passed
--- PASS: TestScriptWithNoContentType (0.00s)
PASS

=== RUN   TestScriptWithNoContentBody
    script_parser_test.go:403: Find inscription with content type: text/plain;charset=utf-8, content length: 0
    script_parser_test.go:409: test script with no content body: test passed
--- PASS: TestScriptWithNoContentBody (0.00s)
PASS

=== RUN   TestScriptWithZeroPush
    script_parser_test.go:438: Find inscription with content type: text/plain;charset=utf-8, content length: 0
    script_parser_test.go:444: test script with zero push: test passed
--- PASS: TestScriptWithZeroPush (0.00s)
PASS

=== RUN   TestScriptWithMultiplePushes
    script_parser_test.go:475: Find inscription with content type: text/plain;charset=utf-8, content length: 66
    script_parser_test.go:481: test script with multiple pushes: test passed
--- PASS: TestScriptWithMultiplePushes (0.00s)
PASS

=== RUN   TestScriptWithNoEndIf
=== PAUSE TestScriptWithNoEndIf
=== CONT  TestScriptWithNoEndIf
    script_parser_test.go:535: test script with no END_IF: test passed
    script_parser_test.go:535: test script with no END_IF: test passed
--- PASS: TestScriptWithNoEndIf (0.00s)
PASS
ok      command-line-arguments  0.173s
```