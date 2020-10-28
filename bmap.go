package bmap

import (
	"github.com/bitcoinschema/go-aip"
	"github.com/bitcoinschema/go-bap"
	"github.com/rohenaz/go-b"
	"github.com/rohenaz/go-bob"
	mapp "github.com/rohenaz/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
	Blk bob.Blk      `json:"blk,omitempty" bson:"blk,omitempty"`
	Tx  bob.TxInfo   `json:"tx,omitempty" bson:"tx,omitempty"`
	In  []bob.Input  `json:"in,omitempty" bson:"in,omitempty"`
	Out []bob.Output `json:"out,omitempty" bson:"out,omitempty"`
	B   *b.B         `json:"B,omnitempty" bson:"B,omitempty"`
	MAP *mapp.MAP    `json:"MAP,omitempty" bson:"MAP,omitempty"`
	AIP *aip.Aip     `json:"AIP,omitempty" bson:"AIP,omitempty"`
	BAP *bap.Data    `json:"BAP,omnitempty" bson:"BAP,omitempty"`
}

// New creates a new BmapTx
func New() *Tx {
	return &Tx{}
}

// FromBob returns a BmapTx from a BobTx
func (bTx *Tx) FromBob(bobTx *bob.BobTx) (err error) {
	for _, out := range bobTx.Out {
		for _, tape := range out.Tape {
			switch tape.Cell[0].S {
			case aip.Prefix:
				bTx.AIP = aip.New()
				bTx.AIP.FromTape(tape)
				bTx.AIP.SetData(out.Tape)
			case bap.Prefix:
				bTx.BAP, err = bap.NewFromTape(&tape)
			case mapp.Prefix:
				bTx.MAP = mapp.New()
				err = bTx.MAP.FromTape(tape)
			case b.Prefix:
				bTx.B = b.New()
				err = bTx.B.FromTape(tape)
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
