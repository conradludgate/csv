package csv

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	Foo    string
	Bar    int64 `csv:"bar"`
	Time   time.Time
	Custom Custom
}

type Custom struct {
	A string
	B int
}

func (c *Custom) UnmarshalCSV(value string) error {
	split := strings.Split(value, "|")
	if len(split) != 2 {
		return fmt.Errorf("invalid data for custom decode")
	}
	c.A = split[0]
	var err error
	c.B, err = strconv.Atoi(split[1])
	return err
}

func TestDecodePass(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	if !assert.Nil(t, err) {
		return
	}

	time1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:05-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2020-07-03T16:39:44+01:00")

	expected := []Data{
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

	assert.Equal(t, expected, v)
}

func TestDecodeFailMissingField(t *testing.T) {
	data := bytes.NewReader([]byte(`Foobar,Time,Custom
hello world,2006-01-02T15:04:05-07:00,value1|1
goodbye world,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "Decode: field for header[Foobar] was not found")
}

func TestDecodeFailDecode(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775808,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"9223372036854775808\": value out of range")
}

func TestDecodeFailDecodeCustom(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "invalid data for custom decode")
}

func TestDecodeFailInconsistentRow(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "record on line 2: wrong number of fields")
}

func TestDecodeFailInconsistentRow2(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1,extra
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "record on line 2: wrong number of fields")
}

func TestDecodePass_DifferentOrder(t *testing.T) {
	data := bytes.NewReader([]byte(`Custom,Foo,bar,Time
value1|1,hello world,9223372036854775807,2006-01-02T15:04:05-07:00
value2|2,goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(&v)
	if !assert.Nil(t, err) {
		return
	}

	time1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:05-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2020-07-03T16:39:44+01:00")

	expected := []Data{
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

	assert.Equal(t, expected, v)
}

func TestDecodeFailNotPointer(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []Data
	err := decoder.Decode(v)
	assert.EqualError(t, err, "Decode: could not decode into type []csv.Data - must be a pointer")
}

func TestDecodeFailNotSlicePointer(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v Data
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "Decode: could not decode into type csv.Data")
}

func TestDecodeFailNotStructSlicePointer(t *testing.T) {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	decoder := NewDecoder(data)
	var v []string
	err := decoder.Decode(&v)
	assert.EqualError(t, err, "Decode: could not decode into type []string - expected a slice of structs")
}
