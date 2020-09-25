package bmap

import (
	"github.com/rohenaz/go-aip"
	"github.com/rohenaz/go-b"
	"github.com/rohenaz/go-bap"
	"github.com/rohenaz/go-bob"
	mapp "github.com/rohenaz/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
	bob.Tx
	AIP *aip.Aip  `json:"AIP,omitempty" bson:"AIP,omitempty"`
	MAP *mapp.MAP `json:"MAP,omitempty" bson:"MAP,omitempty"`
	BAP *bap.Data `json:"BAP,omnitempty" bson:"BAP,omitempty"`
	B   *b.B      `json:"B,omnitempty" bson:"B,omitempty"`
}

// New creates a new BmapTx
func New() *Tx {
	return &Tx{}
}

// FromBob returns a BmapTx from a BobTx
func (bTx *Tx) FromBob(bobTx *bob.Tx) (err error) {
	for _, out := range bobTx.Out {
		for _, tape := range out.Tape {
			switch tape.Cell[0].S {
			case aip.Prefix:
				bTx.AIP = aip.New()
				bTx.AIP.FromTape(tape)
				bTx.AIP.SetData(out.Tape)
			case bap.Prefix:
				bTx.BAP = bap.New()
				err = bTx.BAP.FromTape(tape)
			case mapp.Prefix:
				bTx.MAP = mapp.New()
				err = bTx.MAP.FromTape(tape)
			case b.Prefix:
				bTx.B = b.New()
				err = bTx.B.FromTape(tape)
			}
		}
	}
	return nil
}
