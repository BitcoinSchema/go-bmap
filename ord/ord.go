// Package ord is for parsing 1sat ordinals
package ord

import (
	"encoding/base64"

	"github.com/bitcoinschema/go-bpu"
)

// Prefix is the OP_RETURN prefix for the 1Sat Ordinals inscription protocol
const Prefix string = "ord"

// Ordinal tells wether an inscription is found
type Ordinal struct {
	Data        []byte
	ContentType string
	Vout        uint8
}

// FromTape sets the ordinal data from a bpu.Tape
func (o *Ordinal) FromTape(tape *bpu.Tape) (err error) {
	data := tape.Cell[9].B
	contentType := tape.Cell[11].S
	if data != nil && contentType != nil {
		var dataBytes []byte
		dataBytes, err = base64.StdEncoding.DecodeString(*data)
		if err != nil {
			return
		}
		o.Data = dataBytes
		o.ContentType = *contentType
		o.Vout = tape.I
	}
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
