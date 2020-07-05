package csv

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (c Custom) MarshalCSV() string {
	return fmt.Sprintf("%s|%d", c.A, c.B)
}

func TestEncodePass(t *testing.T) {
	time1 := time.Date(2006, 01, 02, 15, 04, 05, 0, time.FixedZone("MST", -7*60*60))
	time2 := time.Date(2020, 07, 03, 16, 39, 44, 0, time.FixedZone("BST", 1*60*60))

	data := []Data{
		{
			Foo:  "hello world",
			Bar:  9223372036854775807,
			Time: time1,
			Custom: Custom{
				A: "value1",
				B: 1,
			},
		},
		{
			Foo:  "goodbye world",
			Bar:  -9223372036854775808,
			Time: time2,
			Custom: Custom{
				A: "value2",
				B: 2,
			},
		},
	}

	expected := `Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2
`

	bytes, err := Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(bytes))
}

func TestEncodeFailNotSlice(t *testing.T) {
	time1 := time.Date(2006, 01, 02, 15, 04, 05, 0, time.FixedZone("MST", -7*60*60))

	data := Data{
		Foo:  "hello world",
		Bar:  9223372036854775807,
		Time: time1,
		Custom: Custom{
			A: "value1",
			B: 1,
		},
	}

	b, err := Marshal(data)
	assert.EqualError(t, err, "Encode: could not encode type csv.Data")
	assert.Empty(t, b)
}

func TestEncodeFailNotStructSlice(t *testing.T) {
	data := []string{"a", "b"}
	b, err := Marshal(data)
	assert.EqualError(t, err, "Encode: could not encode type []string - expected a collection of structs")
	assert.Empty(t, b)
}

func TestEncodePass_UseCRLF(t *testing.T) {
	time1 := time.Date(2006, 01, 02, 15, 04, 05, 0, time.FixedZone("MST", -7*60*60))
	time2 := time.Date(2020, 07, 03, 16, 39, 44, 0, time.FixedZone("BST", 1*60*60))

	data := []Data{
		{
			Foo:  "hello world",
			Bar:  9223372036854775807,
			Time: time1,
			Custom: Custom{
				A: "value1",
				B: 1,
			},
		},
		{
			Foo:  "goodbye world",
			Bar:  -9223372036854775808,
			Time: time2,
			Custom: Custom{
				A: "value2",
				B: 2,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(buf)
	encoder.UseCRLF()
	err := encoder.Encode(data)
	assert.Nil(t, err)

	expected := "Foo,bar,Time,Custom\r\nhello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1\r\ngoodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2\r\n"

	assert.Equal(t, expected, buf.String())
}

func TestEncodePass_NoCRLF(t *testing.T) {
	time1 := time.Date(2006, 01, 02, 15, 04, 05, 0, time.FixedZone("MST", -7*60*60))
	time2 := time.Date(2020, 07, 03, 16, 39, 44, 0, time.FixedZone("BST", 1*60*60))

	data := []Data{
		{
			Foo:  "hello world",
			Bar:  9223372036854775807,
			Time: time1,
			Custom: Custom{
				A: "value1",
				B: 1,
			},
		},
		{
			Foo:  "goodbye world",
			Bar:  -9223372036854775808,
			Time: time2,
			Custom: Custom{
				A: "value2",
				B: 2,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(buf)
	encoder.NoCRLF()
	err := encoder.Encode(data)
	assert.Nil(t, err)

	expected := "Foo,bar,Time,Custom\nhello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1\ngoodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2\n"

	assert.Equal(t, expected, buf.String())
}
