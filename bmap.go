package bmap

import (
	"github.com/bitcoinschema/go-aip"
	"github.com/bitcoinschema/go-b"
	"github.com/bitcoinschema/go-bap"
	"github.com/bitcoinschema/go-bmap/run"
	"github.com/bitcoinschema/go-bob"
	"github.com/bitcoinschema/go-bpu"
	magic "github.com/bitcoinschema/go-map"
)

// Tx is a Bmap formatted tx// Tx is a Bmap formatted tx
type Tx struct {
	bpu.BpuTx
	AIP *aip.Aip  `json:"AIP,omitempty" bson:"AIP,omitempty"`
	B   *b.B      `json:"B,omitempty" bson:"B,omitempty"`
	BAP *bap.Bap  `json:"BAP,omitempty" bson:"BAP,omitempty"`
	MAP magic.MAP `json:"MAP,omitempty" bson:"MAP,omitempty"`
	Run *run.Jig  `json:"Run,omitempty" bson:"Run,omitempty"`
}

// NewFromBob returns a new BmapTx from a BobTx
func NewFromBob(bobTx *bob.Tx) (bmapTx *Tx, err error) {
	bmapTx = new(Tx)
	err = bmapTx.FromBob(bobTx)
	return
}

// NewFromTx returns a new BmapTx from a hex string
func NewFromTx(tx string) (bmapTx *Tx, err error) {
	var bobTx *bob.Tx
	if bobTx, err = bob.NewFromRawTxString(tx); err != nil {
		return
	}

	bmapTx = new(Tx)
	err = bmapTx.FromBob(bobTx)
	return
}

// FromBob returns a BmapTx from a BobTx
func (t *Tx) FromBob(bobTx *bob.Tx) (err error) {
	for _, out := range bobTx.Out {
		for index, tape := range out.Tape {
			if len(tape.Cell) > 0 && tape.Cell[0].S != nil {
				prefixData := *tape.Cell[0].S
				switch prefixData {
				case run.Prefix:
					if t.Run, err = run.NewFromTape(tape); err != nil {
						return
					}
				case aip.Prefix:
					t.AIP = aip.NewFromTape(tape)
					t.AIP.SetDataFromTapes(out.Tape)
				case bap.Prefix:
					if t.BAP, err = bap.NewFromTape(&out.Tape[index]); err != nil {
						return
					}
				case magic.Prefix:
					if t.MAP, err = magic.NewFromTape(&out.Tape[index]); err != nil {
						return
					}
				case b.Prefix:
					if t.B, err = b.NewFromTape(out.Tape[index]); err != nil {
						return
					}
				}
			}
		}

		// Set inherited fields
		t.Blk = bobTx.Blk
		t.In = bobTx.In
		t.Out = bobTx.Out
		t.Tx = bobTx.Tx
	}
	return
}
