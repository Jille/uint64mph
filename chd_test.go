package uint64mph

import (
	"bytes"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleData = map[uint64]uint64{
		4497751427889084562: 7350820204916009064,
		6500312395473234249: 3838014355565228088,
		4512880937691157593: 8476067645091641406,
		341165985643372816:  9137376927100172348,
		5935824643270476650: 815096534710225439,
		2672920638734362811: 6212145708329429500,
		85228549085877255:   5674464906966263850,
	}
)

var (
	words []uint64
)

func init() {
	words = make([]uint64, 102401)
	for i := range words {
		words[i] = rand.Uint64()
	}
}

func TestCHDBuilder(t *testing.T) {
	b := Builder()
	for k, v := range sampleData {
		b.Add(k, v)
	}
	c, err := b.Build()
	assert.NoError(t, err)
	assert.Equal(t, 7, len(c.keys))
	for k, v := range sampleData {
		assert.Equal(t, v, c.Get(k))
	}
	assert.Equal(t, uint64(math.MaxUint64), c.Get(5))
}

func TestCHDSerialization(t *testing.T) {
	cb := Builder()
	for _, v := range words {
		cb.Add(v, v)
	}
	m, err := cb.Build()
	assert.NoError(t, err)
	w := &bytes.Buffer{}
	err = m.Write(w)
	assert.NoError(t, err)

	n, err := Mmap(w.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, n.r, m.r)
	assert.Equal(t, n.indices, m.indices)
	assert.Equal(t, n.keys, m.keys)
	assert.Equal(t, n.values, m.values)
	for _, v := range words {
		assert.Equal(t, v, n.Get(v))
	}
}

func TestCHDSerialization_empty(t *testing.T) {
	cb := Builder()
	m, err := cb.Build()
	assert.NoError(t, err)
	w := &bytes.Buffer{}
	err = m.Write(w)
	assert.NoError(t, err)

	n, err := Mmap(w.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, n.r, m.r)
	assert.Equal(t, n.indices, m.indices)
	assert.Equal(t, n.keys, m.keys)
	assert.Equal(t, n.values, m.values)
}

func TestCHDSerialization_one(t *testing.T) {
	cb := Builder()
	cb.Add(13, 37)
	m, err := cb.Build()
	assert.NoError(t, err)
	w := &bytes.Buffer{}
	err = m.Write(w)
	assert.NoError(t, err)

	n, err := Mmap(w.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, n.r, m.r)
	assert.Equal(t, n.indices, m.indices)
	assert.Equal(t, n.keys, m.keys)
	assert.Equal(t, n.values, m.values)
}

func BenchmarkBuiltinMap(b *testing.B) {
	keys := []uint64{}
	d := map[uint64]uint64{}
	for _, k := range words {
		d[k] = k
		keys = append(keys, k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d[keys[i%len(keys)]]
	}
}

func BenchmarkCHD(b *testing.B) {
	keys := words
	mph := Builder()
	for _, k := range words {
		mph.Add(k, k)
	}
	h, _ := mph.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Get(keys[i%len(keys)])
	}
}
