package benchmarks_test

import (
	"io/ioutil"
	"log"
	"testing"

	standardModel "github.com/easyaction/encoding-benchmarks/messagepack/model/standard"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v4"
)

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
	entityMsgpack := &standardModel.Superhero{}
	err := msgpack.Unmarshal(jsonPayload, entityMsgpack)
	if err != nil {
		t.Fatal(err)
	}

	err = jsoniter.ConfigFastest.Unmarshal(jsonPayload, entity)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, entity.Id, entityMsgpack.ID)
	assert.Equal(t, entity.AffiliationId, entityMsgpack.AffiliationID)
	assert.Equal(t, entity.Name, entityMsgpack.Name)
	assert.Equal(t, entity.Life, entityMsgpack.Life)
	assert.Equal(t, entity.Energy, entityMsgpack.Energy)

	for i, power := range entity.Powers {
		parserPower := entityMsgpack.Powers[i]
		assert.Equal(t, power.Id, entityMsgpack.ID)
		assert.Equal(t, power.Name, entityMsgpack.Name)
		assert.Equal(t, power.Energy, entityMsgpack.Energy)
		assert.Equal(t, power.Damage, entityMsgpack.Damage)
		assert.Equal(t, power.Passive, entityMsgpack.Passive)
	}
}

// BenchmarkJSONDecodeParser performs benchmark of JSON decoding by using low level API provided by buger/jsonparser library.
func BenchmarkMsgpackDecodeParser(b *testing.B) {
	e := &standardModel.Superhero{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := msgpack.Unmarshal(jsonPayload, e)
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
