package sfp

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/libsv/go-bt"
)

// SFP - Simple Fabriik Protocol
type SFP struct {
	Amount            uint64
	AssetPaymail      string
	AuthorizerAddress string
	OwnerAddress      string
	IssuerAddress     string
	LinkedPrevOut     string
	LinkedPrevOutSig  string
	Version           string
	Notes             string
}

// NewFromUtxo returns a SFP token object from a utxo
func NewFromUtxo(utxo *bt.Output, linkedPrevOut *string, linkedPrevOutSig *string) (sfp *SFP, e error) {

	sfp = &SFP{}

	asm, err := utxo.LockingScript.ToASM()
	if err != nil {
		return nil, err
	}

	scriptParts := strings.Split(asm, " ")

	if len(scriptParts) < 9 || scriptParts[0] != "OP_NOP" {
		return nil, fmt.Errorf("not a valid SFP Token! %d", len(scriptParts))
	}

	// Detect 2nd script chunk is `"SFP@0.3",`
	var protocolBytes []byte
	protocolBytes, err = hex.DecodeString(scriptParts[1])
	if err != nil {
		return nil, err
	}

	if strings.Contains(strings.ToLower(string(protocolBytes)), "sfp@") {
		// Set version
		protocolVersionParts := strings.Split(string(protocolBytes), "@")
		if len(protocolVersionParts) != 2 {
			return nil, fmt.Errorf("invalid version signature: %s", string(protocolBytes))
		}

		sfp.Version = protocolVersionParts[1]
		// Set Asset Paymail
		var assetPaymailBytes []byte
		assetPaymailBytes, err = hex.DecodeString(scriptParts[2])
		if err != nil {
			return nil, err
		}
		sfp.AssetPaymail = string(assetPaymailBytes)

		log.Printf("Paymail %s", sfp.AssetPaymail)
		// Set Authorizer Address
		// var authorizerAddressBytes []byte
		// authorizerAddressBytes, err = hex.DecodeString(scriptParts[3])
		// if err != nil {
		// 	return nil, err
		// }
		sfp.AuthorizerAddress = scriptParts[3]
		log.Printf("Authorizer %s", sfp.AuthorizerAddress)
		// Set Owner Address
		// var ownerAddressBytes []byte
		// ownerAddressBytes, err = hex.DecodeString(scriptParts[4])
		// if err != nil {
		// 	return nil, err
		// }
		sfp.OwnerAddress = scriptParts[4] // string(ownerAddressBytes)
		log.Printf("Owner %s", sfp.OwnerAddress)
		// Set Issuer Address
		// var issuerAddressBytes []byte
		// issuerAddressBytes, err = hex.DecodeString(scriptParts[5])
		// if err != nil {
		// 	return nil, err
		// }
		sfp.IssuerAddress = scriptParts[5] // string(issuerAddressBytes)

		// Find & set Data
		var pos = 0
		var data []string

		for i, op := range scriptParts {
			if op == "OP_RETURN" {
				pos = i
				continue
			}
			if pos > 0 && i > pos {
				data = append(data, op)
			}
		}

		if pos == 0 {
			return nil, fmt.Errorf("Error finding data")
		}

		// Extract amount
		// To extract the quantity of tokens,
		// take the first 8 bytes from the hex encoded data
		// after the OP_RETURN and convert it to number
		// using little endian encoding.
		r := []rune(data[0][0:16])
		x := make([]rune, 16)
		x[14], x[15] = r[0], r[1]
		x[12], x[13] = r[2], r[3]
		x[10], x[11] = r[4], r[5]
		x[8], x[9] = r[6], r[7]
		x[6], x[7] = r[8], r[9]
		x[4], x[5] = r[10], r[11]
		x[2], x[3] = r[12], r[13]
		x[0], x[1] = r[14], r[15]

		amountStr := strings.TrimPrefix(string(x), "0")
		amount, err := strconv.ParseUint(amountStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to get amount %w", err)
		}
		sfp.Amount = amount
		log.Println("Amount", amount)

		// Set notes
		notes := data[0][16 : len(data[0])-6]

		var notesBytes []byte
		notesBytes, err = hex.DecodeString(notes)
		if err != nil {
			return nil, fmt.Errorf("failed to get notes %w", err)
		}
		sfp.Notes = string(notesBytes)
		log.Println("notes", sfp.Notes)

		// Note on outpoints: I think format is <hex>:<out> or OP_FALSE if not present
		if linkedPrevOut == nil {
			sfp.LinkedPrevOut = "OP_FALSE"
		} else {

			sfp.LinkedPrevOut = *linkedPrevOut
		}
		if linkedPrevOutSig == nil {
			sfp.LinkedPrevOutSig = "OP_FALSE"
		} else {
			sfp.LinkedPrevOutSig = *linkedPrevOutSig
		}
		return sfp, nil
	}
	return nil, fmt.Errorf("Not a valid SFP Token! %s %d", protocolBytes, len(scriptParts))

}

// func request() {

// }

// OP_NOP [sfp@0.3] [asset paymail] [authoriser address] [owner address] [issuer address] [linked previous outpoint]
// [linked previous outpoint signature] <more script code> OP_RETURN [data]
// https://docs.moneybutton.com/docs/sfp/wallets-integration-guide.html
