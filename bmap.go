package bmap

import (
	"github.com/rohenaz/go-aip"
	"github.com/rohenaz/go-bap"
	"github.com/rohenaz/go-bob"
	mapp "github.com/rohenaz/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
	AIP *aip.Aip  `json:"AIP,omitempty" bson:"AIP,omitempty"`
	MAP *mapp.MAP `json:"MAP,omitempty" bson:"MAP,omitempty"`
	BAP *bap.Data `json:"BAP,omnitempty" bson:"BAP,omitempty"`
}

// New creates a new BmapTx
func New() *Tx {
	return &Tx{}
}

// FromBob returns a BmapTx from a BobTx
func (b *Tx) FromBob(bobTx *bob.Tx) {
	for _, out := range bobTx.Out {
		for _, tape := range out.Tape {
			switch tape.Cell[0].S {
			case aip.Prefix:
				b.AIP = aip.New()
				b.AIP.FromTape(tape)
			case bap.Prefix:
				b.BAP = bap.New()
				b.BAP.FromTape(tape)
			case mapp.Prefix:
				b.MAP = mapp.New()
				b.MAP.FromTape(tape)
			}
		}
	}
}
