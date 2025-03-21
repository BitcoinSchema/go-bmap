// Package bmap detects known data protocols within Bitcoin transactions
package bmap

import (
	"github.com/bitcoinschema/go-aip"
	"github.com/bitcoinschema/go-b"
	"github.com/bitcoinschema/go-bap"
	"github.com/bitcoinschema/go-bmap/ord"
	"github.com/bitcoinschema/go-bmap/run"
	"github.com/bitcoinschema/go-bob"
	"github.com/bitcoinschema/go-boost"
	"github.com/bitcoinschema/go-bpu"
	magic "github.com/bitcoinschema/go-map"
	"github.com/bitcoinschema/go-sigma"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// Tx is a Bmap formatted tx
type Tx struct {
	bpu.Tx
	AIP   []*aip.Aip     `json:"AIP,omitempty" bson:"AIP,omitempty"`
	B     []*b.B         `json:"B,omitempty" bson:"B,omitempty"`
	BAP   []*bap.Bap     `json:"BAP,omitempty" bson:"BAP,omitempty"`
	BOOST []*boost.Boost `json:"BOOST,omitempty" bson:"BOOST,omitempty"`
	MAP   []magic.MAP    `json:"MAP,omitempty" bson:"MAP,omitempty"`
	Run   []*run.Jig     `json:"Run,omitempty" bson:"Run,omitempty"`
	Ord   []*ord.Ordinal `json:"Ord,omitempty" bson:"Ord,omitempty"`
	Sigma []*sigma.Sig   `json:"Sigma,omitempty" bson:"Sigma,omitempty"`
}

// NewFromBob returns a new BmapTx from a BobTx
func NewFromBob(bobTx *bob.Tx) (bmapTx *Tx, err error) {
	bmapTx = new(Tx)
	err = bmapTx.FromBob(bobTx)
	if err != nil {
		bmapTx = nil
	}
	return
}

// NewFromTx returns a new BmapTx from a *bt.Tx
func NewFromTx(tx *transaction.Transaction) (bmapTx *Tx, err error) {
	var bobTx *bob.Tx
	if bobTx, err = bob.NewFromTx(tx); err != nil {
		return
	}

	bmapTx = new(Tx)
	err = bmapTx.FromBob(bobTx)
	return
}

// NewFromRawTxString returns a new BmapTx from a hex string
func NewFromRawTxString(tx string) (bmapTx *Tx, err error) {
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
	aipCount := 0
	for vout, out := range bobTx.Out {
		for index, tape := range out.Tape {
			// Handle string prefixes
			if len(tape.Cell) > 0 && tape.Cell[0].S != nil {
				prefixData := *tape.Cell[0].S
				switch prefixData {
				case run.Prefix:
					var runOut *run.Jig
					if runOut, err = run.NewFromTape(tape); err != nil {
						return err
					}
					t.Run = append(t.Run, runOut)
					continue
				case aip.Prefix:
					aipOut := aip.NewFromTape(tape)
					aipOut.SetDataFromTapes(out.Tape, aipCount)
					aipCount++
					t.AIP = append(t.AIP, aipOut)
					continue
				case bap.Prefix:
					var bapOut *bap.Bap
					if bapOut, err = bap.NewFromTape(&out.Tape[index]); err != nil {
						return err
					}
					t.BAP = append(t.BAP, bapOut)
					continue
				case magic.Prefix:
					var mapOut magic.MAP
					if mapOut, err = magic.NewFromTape(&out.Tape[index]); err != nil {
						return err
					}
					t.MAP = append(t.MAP, mapOut)
					continue
				case boost.Prefix:
					var boostOut *boost.Boost
					if boostOut, err = boost.NewFromTape(&out.Tape[index]); err != nil {
						return err
					}
					t.BOOST = append(t.BOOST, boostOut)
					continue
				case b.Prefix:
					var bOut *b.B
					if bOut, err = b.NewFromTape(out.Tape[index]); err != nil {
						return err
					}
					t.B = append(t.B, bOut)
					continue
				case sigma.Prefix:
					sigOut := sigma.NewSigFromTape(out.Tape[index], vout)
					t.Sigma = append(t.Sigma, sigOut)
					continue
				}
			}
			// Handle OPCODE prefixes
			if len(tape.Cell) > 5 && tape.Cell[0].Ops != nil {
				ordScript := ord.ScriptFromTape(tape)
				if len(ordScript) > 0 {

					var ordOut *ord.Ordinal
					if ordOut, err = ord.NewFromTape(out.Tape[index]); err != nil {
						return err
					}
					ordOut.Vout = uint8(len(t.Ord))
					t.Ord = append(t.Ord, ordOut)
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
