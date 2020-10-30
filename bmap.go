package bmap

import (
	"github.com/bitcoinschema/go-aip"
	"github.com/bitcoinschema/go-b"
	"github.com/bitcoinschema/go-bap"
	"github.com/bitcoinschema/go-bob"
	magic "github.com/bitcoinschema/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
	Blk bob.Blk      `json:"blk,omitempty" bson:"blk,omitempty"`
	Tx  bob.TxInfo   `json:"tx,omitempty" bson:"tx,omitempty"`
	In  []bob.Input  `json:"in,omitempty" bson:"in,omitempty"`
	Out []bob.Output `json:"out,omitempty" bson:"out,omitempty"`
	B   *b.B         `json:"B,omnitempty" bson:"B,omitempty"`
	MAP magic.MAP    `json:"MAP,omitempty" bson:"MAP,omitempty"`
	AIP *aip.Aip     `json:"AIP,omitempty" bson:"AIP,omitempty"`
	BAP *bap.Bap     `json:"BAP,omnitempty" bson:"BAP,omitempty"`
}

// New creates a new BmapTx
func New() *Tx {
	return &Tx{}
}

// NewFromBob returns a new BmapTx from a BobTx
func NewFromBob(bobTx *bob.Tx) (bmapTx *Tx, err error) {
	bmapTx = New()
	err = bmapTx.FromBob(bobTx)
	return
}

// FromBob returns a BmapTx from a BobTx
func (bTx *Tx) FromBob(bobTx *bob.Tx) (err error) {
	for _, out := range bobTx.Out {
		for _, tape := range out.Tape {
			if len(tape.Cell) > 0 {
				prefixData := tape.Cell[0].S
				switch prefixData {
				case aip.Prefix:
					bTx.AIP = aip.NewFromTape(tape)
					bTx.AIP.SetDataFromTapes(out.Tape)
				case bap.Prefix:
					bTx.BAP, err = bap.NewFromTape(&tape)
				case magic.Prefix:
					bTx.MAP, err = magic.NewFromTape(&tape)
				case b.Prefix:
					bTx.B = b.New()
					err = bTx.B.FromTape(tape)
				}
			}
		}

		// Set inherited fields
		bTx.Tx = bobTx.Tx
		bTx.Blk = bobTx.Blk
		bTx.In = bobTx.In
		bTx.Out = bobTx.Out
	}
	return nil
}
