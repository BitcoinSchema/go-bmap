package bmap

import (
	"log"

	"github.com/rohenaz/go-bap"
	"github.com/rohenaz/go-bob"
	mapp "github.com/rohenaz/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
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
			case bap.BapPrefix:
				log.Println(tape.Cell[0].S)
				// Detect the type
				b.BAP = bap.New()
				b.BAP.FromTape(tape)
				break
			case mapp.MapPrefix:
				b.MAP = mapp.New()
				b.MAP.FromTape(tape)
				break
			}
		}
	}
}
