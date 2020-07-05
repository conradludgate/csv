package csv

import (
	"bytes"
	"errors"
	"math"
	"testing"
	"time"

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

type BadData struct {
	Column NoMarshal
}

type NoMarshal struct {
	A string
}

func TestFailNoMarshal(t *testing.T) {
	data := []BadData{
		{
			Column: NoMarshal{
				A: "fail",
			},
		},
	}

	b, err := Marshal(data)
	assert.EqualError(t, err, "Encode: csv.NoMarshal is not a valid field type - try implement MarshalCSV for it")
	assert.Empty(t, b)
}

func TestFailNoUnmarshal(t *testing.T) {
	data := "A\nfail"

	output := []BadData{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "Decode: csv.NoMarshal is not a valid field type - try implement UnmarshalCSV for it")
	assert.Empty(t, output)
}

type FailReader struct{}

func (r *FailReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("error reading")
}

func TestFailRead(t *testing.T) {
	reader := &FailReader{}
	decoder := NewDecoder(reader)
	output := []BuiltInTypes{}
	err := decoder.Decode(&output)
	assert.EqualError(t, err, "error reading")
	assert.Empty(t, output)
}

type FailWriter struct{}

func (r *FailWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("error writing")
}

func TestFailWrite(t *testing.T) {
	writer := &FailWriter{}
	encoder := NewEncoder(writer)
	input := []BuiltInTypes{}
	err := encoder.Encode(input)
	assert.EqualError(t, err, "error writing")
}

func TestFailWrite2(t *testing.T) {
	writer := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(writer)
	encoder.SetDelimiter(0)
	input := []BuiltInTypes{}
	err := encoder.Encode(input)
	assert.EqualError(t, err, "csv: invalid field or comment delimiter")
}

func TestFailWrite3(t *testing.T) {
	writer := &FailWriter{}
	encoder := NewEncoder(writer)
	input := []BuiltInTypes{}

	// Make lots of data so the underlying writer buffer is forced to flush automatically
	for i := 0; i < 100; i++ {
		input = append(input, BuiltInTypes{
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
		})
	}

	err := encoder.Encode(input)
	assert.EqualError(t, err, "error writing")
}

func TestDecodeFailInt(t *testing.T) {
	type Data struct {
		A int
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"a\": invalid syntax")
}

func TestDecodeFailInt8(t *testing.T) {
	type Data struct {
		A int8
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"a\": invalid syntax")
}

func TestDecodeFailInt16(t *testing.T) {
	type Data struct {
		A int16
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"a\": invalid syntax")
}

func TestDecodeFailInt32(t *testing.T) {
	type Data struct {
		A int32
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"a\": invalid syntax")
}

func TestDecodeFailInt64(t *testing.T) {
	type Data struct {
		A int64
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"a\": invalid syntax")
}

func TestDecodeFailUint(t *testing.T) {
	type Data struct {
		A uint
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestDecodeFailUint8(t *testing.T) {
	type Data struct {
		A uint8
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestDecodeFailUint16(t *testing.T) {
	type Data struct {
		A uint16
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestDecodeFailUint32(t *testing.T) {
	type Data struct {
		A uint32
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestDecodeFailUint64(t *testing.T) {
	type Data struct {
		A uint64
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestDecodeFailBool(t *testing.T) {
	type Data struct {
		A bool
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseBool: parsing \"a\": invalid syntax")
}

func TestDecodeFailFloat32(t *testing.T) {
	type Data struct {
		A float32
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseFloat: parsing \"a\": invalid syntax")
}

func TestDecodeFailFloat64(t *testing.T) {
	type Data struct {
		A float64
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "strconv.ParseFloat: parsing \"a\": invalid syntax")
}

func TestDecodeFailTime(t *testing.T) {
	type Data struct {
		A time.Time
	}
	data := "A\na"

	output := []Data{}
	err := Unmarshal([]byte(data), &output)
	assert.EqualError(t, err, "parsing time \"a\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"a\" as \"2006\"")
}
