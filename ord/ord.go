// Package ord is for parsing 1sat ordinals
package ord

import (
	"encoding/base64"
	"log"

	"github.com/bitcoinschema/go-bpu"
	"github.com/libsv/go-bt/v2/bscript"
)

// Prefix is the OP_RETURN prefix for the 1Sat Ordinals inscription protocol
const Prefix string = "ord"

// Ordinal tells wether an inscription is found
type Ordinal struct {
	Data        []byte
	ContentType string
}

// FromTape sets the ordinal data from a bpu.Tape
func (o *Ordinal) FromTape(tape *bpu.Tape) (err error) {

	ordScript := OrdScriptFromTape(*tape)
	minOrdScriptPushes := 7
	if len(ordScript) == minOrdScriptPushes {
		prefix := ordScript[2].S
		if prefix != nil && *prefix == "ord" {

			for idx, push := range ordScript {
				if push.Op != nil && *push.Op == bscript.Op1 {
					if ordScript[idx+1].S != nil {
						o.ContentType = *ordScript[idx+1].S
					}
				}
				if idx > 0 && push.Op != nil && *push.Op == bscript.Op0 && ordScript[idx+1].B != nil {
					data, err := base64.StdEncoding.DecodeString(*ordScript[idx+1].B)
					if err != nil {
						log.Fatal("error:", err)
					}

					o.Data = data

				}
			}
		}
	}

	// data := tape.Cell[9].B
	// contentType := tape.Cell[11].S
	// if data != nil && contentType != nil {
	// 	var dataBytes []byte
	// 	dataBytes, err = base64.StdEncoding.DecodeString(*data)
	// 	if err != nil {
	// 		return
	// 	}
	// 	o.Data = dataBytes
	// 	o.ContentType = *contentType
	// }
	return
}

// NewFromTape will create a new Ord object from a bpu.Tape
func NewFromTape(tape bpu.Tape) (o *Ordinal, e error) {
	o = new(Ordinal)
	err := o.FromTape(&tape)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// OrdScriptFromTape finds the script: OP_0 OP_IF ... OP_ENDIF
func OrdScriptFromTape(tape bpu.Tape) (ordScript []bpu.Cell) {

	startIdx := 0
	endIdx := 0
	// Find OP_IF and OP_ENDIF indexes
	for idx, c := range tape.Cell {
		if idx > 0 && c.Ops != nil && *c.Ops == "OP_IF" && *tape.Cell[idx-1].Op == 0 {
			startIdx = idx - 1
		}

		if startIdx > 0 && c.Ops != nil && *c.Ops == "OP_ENDIF" {
			endIdx = idx
		}
	}
	return tape.Cell[startIdx:endIdx]
}
