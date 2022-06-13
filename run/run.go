package run

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bitcoinschema/go-bob"
	"github.com/libsv/go-bt"
)

var debug = os.Getenv("BMAP_DEBUG") == "1"

// Prefix is the run protocol prefix found in the 1st pushdata
const Prefix string = "run"

// Command is the run command
type Command string

// Commands enum
const (
	DEPLOY  Command = "DEPLOY"
	NEW     Command = "NEW"
	CALL    Command = "CALL"
	UPGRADE Command = "UPGRADE"
)

// Statement is command + data for execution
// Data depends on the command
// DEPLOY	Upload new code	[<src1>, <props1>, <src2>, <props2>, ...]
// NEW	Instantiate a jig	[<jig class>, <args>]
// CALL	Call a method on a jig	[<jig>, <method>, <args>]
// UPGRADE	Replace code with new code	[<code>, <src>, <props>]
type Statement struct {
	Op   Command
	Data []interface{} // Depends on command
}

// Payload contains the following run metadata:
// in	Number of jig and code inputs
// ref	Array of references to jigs and code used by not spent
// out	State hashes of jigs and code in transaction outputs
// del	State hashes of jigs and code deleted
// cre	New owners of jigs and code created
// exec	Statements to execute on the jigs
type Payload struct {
	In   uint64      `json:"in"`
	Ref  []string    `json:"ref"`
	Out  []string    `json:"out"`
	Del  []string    `json:"del"`
	Cre  []string    `json:"cre"`
	Exec []Statement `json:"exec"`
}

// Jig is a RunOnBitcoin object
type Jig struct {
	AppID   string
	Version uint64
	Payload Payload // not sure what data format is actually best for this
}

// NewFromUtxo returns a Jig from a bt.Output
func NewFromUtxo(utxo *bt.Output) (jig *Jig, e error) {

	jig = &Jig{}

	script, err := utxo.LockingScript.ToASM()
	if err != nil {
		return nil, err
	}

	scriptParts := strings.Split(script, " ")

	// Collect OP_RETURN data from script
	var pos = 0
	var data []string

	for i, op := range scriptParts {
		// Find OP_RETURN
		if op == "OP_RETURN" {
			// Turn on collector
			pos = i
			continue
		}
		// Collect data
		if pos > 0 && i > pos {
			data = append(data, op)
		}
	}

	for i, val := range data {
		// Run pushdata format:
		// 0 - run
		// 1- version
		// 2 - App ID
		// 3 - json payload

		switch i {
		case 0:
			prefix, err := hex.DecodeString(val)
			if string(prefix) != Prefix {
				return nil, fmt.Errorf("Not a valid run Tx: %w", err)
			}
		case 1:
			// TODO: Convert from asm to int
			if debug {
				log.Println("Version", val)
			}
			jig.Version = 0 // val
		case 2:
			appID, err := hex.DecodeString(val)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode app id: %w", err)
			}
			jig.AppID = string(appID)
		case 3:
			payloadBytes, err := hex.DecodeString(val)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode payload: %w", err)
			}
			var payload Payload

			err = json.Unmarshal(payloadBytes, &payload)
			if err != nil {
				return nil, err
			}
			jig.Payload = payload
			// jig.Payload = ?
		}
	}

	return jig, nil
}

// NewFromTape will create a new AIP object from a bob.Tape
// Using the FromTape() alone will prevent validation (data is needed via SetData to enable)
func NewFromTape(tape bob.Tape) (j *Jig, e error) {
	j = new(Jig)
	err := j.FromTape(&tape)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// IsToken checks this is a NEW Jig and extends the Token class
// TODO: Make check more strict, occurs befopre first '{' character
func (j *Jig) IsToken() bool {
	// Check that this is a `Class?` and contains `Extends Token`

	for _, cmd := range j.Payload.Exec {
		if cmd.Op == NEW {
			return strings.Contains(fmt.Sprintf("%v", cmd.Data[0]), " Extends Token {") // Jig class string
		}
	}
	return false
}

// FromTape sets Jig data from Bob Tape
func (j *Jig) FromTape(tape *bob.Tape) error {

	// Run pushdata format:
	// 0 - run
	// 1 - version
	// 2 - App ID
	// 3 - json payload

	if len(tape.Cell) == 4 {
		if tape.Cell[0].S != "run" {
			return fmt.Errorf("Not a run tape")
		}

		// Set the APP ID
		// TODO APP ID is not set on most run transactions - just OP_FALSE
		j.AppID = tape.Cell[2].S

		// Set the version
		// bob parses this in a weird way, it should be just a number, but we can only get the OP_DATA_ hex value
		version := strings.Replace(tape.Cell[1].H, "OP_DATA_", "", 1)
		num, err := strconv.ParseInt(version, 16, 64)
		if err != nil {
			return err
		}

		j.Version = uint64(num)

		var payload Payload

		err = json.Unmarshal([]byte(tape.Cell[3].S), &payload)
		if err != nil {
			return err
		}
		j.Payload = payload
	} else {
		return fmt.Errorf("Pushdata length is incorrect. Got %d expected 4", len(tape.Cell))
	}
	return nil
}
