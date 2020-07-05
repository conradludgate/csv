package csv

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BuiltInTypes struct {
	A string
	B int64
	C int32
	D int16
	E int8
	F int
	G uint64
	H uint32
	I uint16
	J uint8
	K uint
	L bool
	M float64
	N float32
}

func TestBuiltInTypes(t *testing.T) {
	input := []BuiltInTypes{
		{
			A: "a",
			B: -9223372036854775808,
			C: -2147483648,
			D: -32768,
			E: -128,
			F: -1,
			G: 18446744073709551615,
			H: 4294967295,
			I: 65535,
			J: 255,
			K: 1,
			L: true,
			M: math.Pi,
			N: math.E,
		},
	}

	output1, err := Marshal(input)
	assert.Nil(t, err)

	expectedOutput1 := `A,B,C,D,E,F,G,H,I,J,K,L,M,N
a,-9223372036854775808,-2147483648,-32768,-128,-1,18446744073709551615,4294967295,65535,255,1,true,3.141592653589793,2.718282
`

	assert.Equal(t, []byte(expectedOutput1), output1)

	output2 := []BuiltInTypes{}
	err = Unmarshal(output1, &output2)
	assert.Nil(t, err)

	expectedOutput2 := []BuiltInTypes{
		{
			A: "a",
			B: -9223372036854775808,
			C: -2147483648,
			D: -32768,
			E: -128,
			F: -1,
			G: 18446744073709551615,
			H: 4294967295,
			I: 65535,
			J: 255,
			K: 1,
			L: true,
			M: 3.141592653589793,
			N: 2.718282,
		},
	}

	assert.Equal(t, expectedOutput2, output2)
}
