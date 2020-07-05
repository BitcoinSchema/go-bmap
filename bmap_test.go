package bmap

import (
	"testing"

	"github.com/rohenaz/go-bob"
	mapp "github.com/rohenaz/go-map"
)

func TestMap(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: mapp.Prefix},
			{S: mapp.SET},
			{S: "keyName1"},
			{S: "something"},
			{S: "keyName2"},
			{S: "something else"},
		},
	}
	m := mapp.MAP{}
	m.FromTape(tape)
	if m["keyName1"] != "something" {
		t.Errorf("SET Failed %s", m["keyName1"])
	}
}
