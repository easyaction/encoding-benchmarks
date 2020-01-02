package benchmarks_test

import (
	"io/ioutil"
	"log"
	"testing"

	standardModel "github.com/easyaction/encoding-benchmarks/messagepack/model/standard"
	vmihailencoModel "github.com/easyaction/encoding-benchmarks/messagepack/model/vmihailenco"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v4"
)

// jsonPayload is a variable holding encoded JSON reference payload used in all benchmarks.
var jsonPayload []byte

// jsonResult is a dummy output variable for each benchmark. In benchmarks all results must be copied over to an exported variable to prevent Go compiler from skipping parts of code which results are never used.
var jsonResult interface{}

// init reads JSON reference payload.
func init() {
	var err error

	jsonPayload, err = ioutil.ReadFile("./messagepack/payload/superhero.json")
	if err != nil {
		log.Fatal(err)
	}
}

// TestJSONParserDecoding tests decoding through low level API provided by buger/jsonparser library,
// Test is required to make sure that custom unmarshal method is working properly.
func TestVmihailencoDecoding(t *testing.T) {
	entityParser := &vmihailencoModel.Superhero{}
	err := entityParser.UnmarshalJSON(jsonPayload)
	if err != nil {
		t.Fatal(err)
	}

	entity := &standardModel.Superhero{}
	err = jsoniter.ConfigFastest.Unmarshal(jsonPayload, entity)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, entity.Id, entityParser.ID)
	assert.Equal(t, entity.AffiliationId, entityParser.AffiliationID)
	assert.Equal(t, entity.Name, entityParser.Name)
	assert.Equal(t, entity.Life, entityParser.Life)
	assert.Equal(t, entity.Energy, entityParser.Energy)

	for i, power := range entity.Powers {
		parserPower := entityParser.Powers[i]
		assert.Equal(t, power.Id, parserPower.ID)
		assert.Equal(t, power.Name, parserPower.Name)
		assert.Equal(t, power.Energy, parserPower.Energy)
		assert.Equal(t, power.Damage, parserPower.Damage)
		assert.Equal(t, power.Passive, parserPower.Passive)
	}
}

// BenchmarkJSONDecodeParser performs benchmark of JSON decoding by using low level API provided by buger/jsonparser library.
func BenchmarkMsgpackDecodeParser(b *testing.B) {
	e := &vmihailencoModel.Superhero{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := e.UnmarshalJSON(jsonPayload)
		if err != nil {
			b.Fatal(err)
		}

		jsonResult = e
	}
}

// BenchmarkJSONEncodeIterator performs benchmark of JSON encoding by json-iterator/go library.
func BenchmarkMsgpackEncodeIterator(b *testing.B) {
	e := &standardModel.Superhero{}
	err := jsoniter.ConfigFastest.Unmarshal(jsonPayload, e)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		p, err := msgpack.Marshal(e)
		if err != nil {
			b.Fatal(err)
		}
		jsonResult = p
	}
}
