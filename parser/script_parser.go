package parser

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	log "github.com/sirupsen/logrus"
)

const (
	ProtocolID     string = "6f7264"
	BodyTag        string = "00"
	ContentTypeTag string = "01"
)

type TransactionInscription struct {
	Inscription *InscriptionContent
	TxInIndex   uint32
	TxInOffset  uint64
}

type InscriptionContent struct {
	ContentType             string
	ContentBody             []byte
	ContentLength           uint64
	IsUnrecognizedEvenField bool
}

func ParseInscriptionsFromTransaction(msgTx *wire.MsgTx) []*TransactionInscription {
	var inscriptionsFromTx []*TransactionInscription
	txHash := msgTx.TxHash().String()

	if !msgTx.HasWitness() {
		log.Debugf("Tx: %s inputs does not contain witness data", txHash)
		return nil
	}

	for i, v := range msgTx.TxIn {
		index, input := i, v
		if len(input.Witness) <= 1 {
			log.Debugf("Tx: %s, the length of tx input witness data is %d", txHash, len(input.Witness))
			continue
		}
		if len(input.Witness) == 2 && input.Witness[len(input.Witness)-1][0] == txscript.TaprootAnnexTag {
			log.Debugf("Tx: %s, tx witness contains Taproot Annex data but the length of tx input witness data is 2",
				txHash)
			continue
		}

		// If Taproot Annex data exists, take the last element of the witness as the script data, otherwise,
		// take the penultimate element of the witness as the script data
		var witnessScript []byte
		if input.Witness[len(input.Witness)-1][0] == txscript.TaprootAnnexTag {
			witnessScript = input.Witness[len(input.Witness)-1]
		} else {
			witnessScript = input.Witness[len(input.Witness)-2]
		}

		// Parse script and get ordinals content
		inscriptions := ParseInscriptions(witnessScript)
		if len(inscriptions) == 0 {
			continue
		}
		for i, v := range inscriptions {
			txInOffset, inscription := i, v
			inscriptionsFromTx = append(inscriptionsFromTx, &TransactionInscription{
				Inscription: inscription,
				TxInIndex:   uint32(index),
				TxInOffset:  uint64(txInOffset),
			})
		}
	}
	return inscriptionsFromTx
}

func ParseInscriptions(witnessScript []byte) []*InscriptionContent {
	var (
		inscriptions []*InscriptionContent
	)

	// Parse inscription content from witness script
	tokenizer := txscript.MakeScriptTokenizer(0, witnessScript)
	for tokenizer.Next() {
		// Check inscription envelop header: OP_FALSE(0x00), OP_IF(0x63), PROTOCOL_ID([0x6f, 0x72, 0x64])
		if tokenizer.Opcode() == txscript.OP_FALSE {
			if !tokenizer.Next() || tokenizer.Opcode() != txscript.OP_IF {
				return inscriptions
			}
			if !tokenizer.Next() || hex.EncodeToString(tokenizer.Data()) != ProtocolID {
				return inscriptions
			}
			inscription := parseOneInscription(&tokenizer)
			if inscription != nil {
				inscriptions = append(inscriptions, inscription)
			}
		}
	}

	return inscriptions
}

func parseOneInscription(tokenizer *txscript.ScriptTokenizer) *InscriptionContent {
	var (
		tags                    = make(map[string][]byte)
		contentType             string
		contentBody             []byte
		contentLength           uint64
		isUnrecognizedEvenField bool
	)

	// Find any pushed data in the script. This includes OP_0, but not OP_1 - OP_16.
	for tokenizer.Next() {
		if tokenizer.Opcode() == txscript.OP_ENDIF {
			break
		} else if hex.EncodeToString([]byte{tokenizer.Opcode()}) == BodyTag {
			var body []byte
			for tokenizer.Next() {
				if tokenizer.Opcode() == txscript.OP_ENDIF {
					break
				} else if tokenizer.Opcode() == txscript.OP_0 {
					// OP_0 push no data
					continue
				} else if tokenizer.Opcode() >= txscript.OP_DATA_1 && tokenizer.Opcode() <= txscript.OP_PUSHDATA4 {
					// Taproot's restriction, individual data pushes may not be larger than 520 bytes.
					if len(tokenizer.Data()) > 520 {
						log.Errorf("data is longer than 520")
						return nil
					}
					body = append(body, tokenizer.Data()...)
				} else {
					// Invalid opcode found in content body, e.g., 615a7c90df1d4fdd07c6ea98766bc6846dd5264a9fa81ca41611bbf9bde38cf8.
					return nil
				}
			}
			tags[BodyTag] = body
			break
		} else {
			if tokenizer.Data() == nil {
				return nil
			}
			tag := hex.EncodeToString(tokenizer.Data())
			if _, ok := tags[tag]; ok {
				return nil
			}
			if tokenizer.Next() {
				tags[tag] = tokenizer.Data()
			}
		}
	}

	// No OP_ENDIF
	if tokenizer.Opcode() != txscript.OP_ENDIF {
		return nil
	}

	// Error occurred
	if err := tokenizer.Err(); err != nil {
		return nil
	}

	// Get inscription content
	for k := range tags {
		key := k
		if key == ContentTypeTag {
			contentType = string(tags[ContentTypeTag])
			continue
		}
		if key == BodyTag {
			contentBody = tags[BodyTag]
			contentLength = uint64(len(contentBody))
			continue
		}
		// Unrecognized even tag
		tag, _ := hex.DecodeString(key)
		if len(tag) > 0 && int(tag[0])%2 == 0 {
			isUnrecognizedEvenField = true
			continue
		}
	}

	inscription := &InscriptionContent{
		ContentType:             contentType,
		ContentBody:             contentBody,
		ContentLength:           contentLength,
		IsUnrecognizedEvenField: isUnrecognizedEvenField,
	}
	return inscription
}
