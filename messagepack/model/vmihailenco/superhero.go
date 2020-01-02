package model

import (
	"github.com/vmihailenco/msgpack/v4"
)

// superheroPaths is a list of paths to be parsed by buger/jsonparser.
var superheroPaths = [][]string{
	{"id"},
	{"affiliation_id"},
	{"name"},
	{"life"},
	{"energy"},
	{"powers"},
}

// Superhero is a model to which a custom unmarshal method utilizing buger/jsonparser API is attached.
type Superhero struct {
	ID            int
	AffiliationID int
	Name          string
	Life          float64
	Energy        float64
	Powers        []*Superpower
}

// UnmarshalJSON is a method implementing unmarshaler interface and utilizing buger/jsonparser low level API.
func (s *Superhero) UnmarshalJSON(b []byte) error {
	return msgpack.Unmarshal(b, s)
}

// superpowerPaths is a list of paths to be parsed by buger/jsonparser.
var superpowerPaths = [][]string{
	{"id"},
	{"name"},
	{"damage"},
	{"energy"},
	{"passive"},
}

// Superpower is a model to which a custom unmarshal method utilizing buger/jsonparser API is attached.
type Superpower struct {
	ID      int
	Name    string
	Damage  float64
	Energy  float64
	Passive bool
}

// UnmarshalJSON is a method implementing unmarshaler interface and utilizing buger/jsonparser low level API.
func (s *Superpower) UnmarshalJSON(b []byte) error {
	return msgpack.Unmarshal(b, s)
}
