package stream

import (
	"github.com/jpfourny/papaya/internal/assert"
	"math/rand"
	"testing"
)

func TestGenerate(t *testing.T) {
	s := Generate(func() int { return 1 })
	out := CollectSlice(Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []int{1, 1, 1, 1, 1})
}

func TestRepeat(t *testing.T) {
	s := Repeat(1)
	out := CollectSlice(Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []int{1, 1, 1, 1, 1})
}

func TestRepeatN(t *testing.T) {
	s := RepeatN(1, 3)
	out := CollectSlice(s)
	assert.ElementsMatch(t, out, []int{1, 1, 1})
}

func TestRandomInt(t *testing.T) {
	s := RandomInt(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[int](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []int{8717895732742165505, 2259404117704393152, 6050128673802995827, 501233450539197794, 3390393562759376202})
}

func TestRandomIntn(t *testing.T) {
	s := RandomIntn(rand.NewSource(0), 100) // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[int](Limit(s, 5))   // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []int{74, 14, 53, 6, 15})
}

func TestRandomUint32(t *testing.T) {
	s := RandomUint32(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[uint32](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []uint32{0xf1f85ff5, 0x3eb60825, 0xa7ecbff2, 0xde97a55, 0x5e1a31f6})
}

func TestRandomUint64(t *testing.T) {
	s := RandomUint64(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[uint64](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []uint64{0x78fc2ffac2fd9401, 0x1f5b0412ffd341c0, 0x53f65ff94f6ec873, 0x86f4bd2ae8eea562, 0xaf0d18fb750b2d4a})
}

func TestRandomFloat32(t *testing.T) {
	s := RandomFloat32(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[float32](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []float32{0.94519615, 0.24496509, 0.65595627, 0.05434384, 0.3675872})
}

func TestRandomFloat64(t *testing.T) {
	s := RandomFloat64(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[float64](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []float64{0.9451961492941164, 0.24496508529377975, 0.6559562651954052, 0.05434383959970039, 0.367587206632458534})
}

func TestRandomNormFloat64(t *testing.T) {
	s := RandomNormFloat64(rand.NewSource(0)) // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[float64](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []float64{-0.28158587086436215, 0.570933095808067, -1.6920196326157044, 0.1996229111693099, 1.9195199291234621})
}

func TestRandomExpFloat64(t *testing.T) {
	s := RandomExpFloat64(rand.NewSource(0))  // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[float64](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []float64{4.668112973579268, 0.1601593871172866, 3.0465834105636, 0.06385839451671879, 1.8578917487258961})
}

func TestRandomBool(t *testing.T) {
	s := RandomBool(rand.NewSource(0))     // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[bool](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []bool{true, true, false, true, false})
}

func TestRandomBytes(t *testing.T) {
	s := RandomBytes(rand.NewSource(0))    // Seed with 0 for testing - should produce the same sequence every time.
	out := CollectSlice[byte](Limit(s, 5)) // Limit to 5 elements for testing.
	assert.ElementsMatch(t, out, []byte{0xfa, 0x12, 0xf9, 0x2a, 0xfb})
}
