package bmap

import (
	"testing"

	"github.com/rohenaz/go-b"
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

func TestB(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: b.Prefix},
			{S: "Hello world"},
			{S: "text/plain"},
			{S: "utf8"},
		},
	}
	m := b.B{}
	m.FromTape(tape)
	if m.Data.UTF8 != "Hello world" {
		t.Errorf("Unexpected data %s", m.Data.UTF8)
	}
}
