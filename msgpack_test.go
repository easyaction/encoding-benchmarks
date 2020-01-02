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

// jsonPayload is a variable holding encoded JSON reference payload used in all benchmarks.
var msgPayload []byte

// jsonResult is a dummy output variable for each benchmark. In benchmarks all results must be copied over to an exported variable to prevent Go compiler from skipping parts of code which results are never used.
var msgResult interface{}

// init reads JSON reference payload.
func init() {
	var err error

	jsonPayload, err := ioutil.ReadFile("./messagepack/payload/superhero.json")
	if err != nil {
		log.Fatal(err)
	}

	sh := &standardModel.Superhero{}
	err = jsoniter.ConfigFastest.Unmarshal(jsonPayload, sh)
	if err != nil {
		log.Fatal(err)
	}

	b, err := msgpack.Marshal(sh)

	err = ioutil.WriteFile("./messagepack/payload/superhero.msgpack", b, 0)
	if err != nil {
		log.Fatal(err)
	}

	msgPayload, err = ioutil.ReadFile("./messagepack/payload/superhero.msgpack")
	if err != nil {
		log.Fatal(err)
	}
}

// TestJSONParserDecoding tests decoding through low level API provided by buger/jsonparser library,
// Test is required to make sure that custom unmarshal method is working properly.
func TestVmihailencoDecoding(t *testing.T) {
	entityMsgpack := &standardModel.Superhero{}
	err := msgpack.Unmarshal(msgPayload, entityMsgpack)
	if err != nil {
		t.Fatal(err)
	}
	entity := &standardModel.Superhero{}

	err = jsoniter.ConfigFastest.Unmarshal(jsonPayload, entity)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, entity.Id, entityMsgpack.Id)
	assert.Equal(t, entity.AffiliationId, entityMsgpack.AffiliationId)
	assert.Equal(t, entity.Name, entityMsgpack.Name)
	assert.Equal(t, entity.Life, entityMsgpack.Life)
	assert.Equal(t, entity.Energy, entityMsgpack.Energy)

	for i, power := range entity.Powers {
		msgpackPower := entityMsgpack.Powers[i]
		assert.Equal(t, power.Id, msgpackPower.Id)
		assert.Equal(t, power.Name, msgpackPower.Name)
		assert.Equal(t, power.Energy, msgpackPower.Energy)
		assert.Equal(t, power.Damage, msgpackPower.Damage)
		assert.Equal(t, power.Passive, msgpackPower.Passive)
	}
}

// BenchmarkJSONDecodeParser performs benchmark of JSON decoding by using low level API provided by buger/jsonparser library.
func BenchmarkMsgpackDecodeParser(b *testing.B) {
	e := &standardModel.Superhero{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := msgpack.Unmarshal(msgPayload, e)
		if err != nil {
			b.Fatal(err)
		}

		msgResult = e
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
		msgResult = p
	}
}
