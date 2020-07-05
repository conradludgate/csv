# csv
A CSV marshaling and unmarshaling library for Go, with a similar API to [encoding/json](https://pkg.go.dev/encoding/json)

[![codecov](https://img.shields.io/codecov/c/gh/conradludgate/csv?style=flat-square)](https://codecov.io/gh/conradludgate/csv)
[![docs](https://img.shields.io/github/v/tag/conradludgate/csv?label=docs&style=flat-square)](https://pkg.go.dev/github.com/conradludgate/csv?tab=doc)

## Usage

To install, run `go get github.com/conradludgate/csv`

### Unmarshal

Takes in CSV data with a header row and will unmarshal it into a slice of a struct that matches the headers

Example:
```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/conradludgate/csv"
)

// Data is the column structure for the CSV.
// Most built in types have a sensible decode function from string
// Custom header names can be specified using the `csv` tag
// Custom types are supported if they implement `csv.UnmarshalCSV`
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

// Example custom unmarshal function
// Turns "A|B" into Custom{A, B}
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

func example1() {
	// Input data with the first row being a header row
	data := []byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	// Create a slice of the column data
	output := []Data{}
	// Unmarshal just like with `json.Unmarshal`
	err := csv.Unmarshal(data, &output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Outputs:
	// [
	//     {
	//         "Foo": "hello world",
	//         "Bar": 9223372036854775807,
	//         "Time": "2006-01-02T15:04:05-07:00",
	//         "Custom": {
	//             "A": "value1",
	//             "B": 1
	//         }
	//     },
	//     {
	//         "Foo": "goodbye world",
	//         "Bar": -9223372036854775808,
	//         "Time": "2020-07-03T16:39:44+01:00",
	//         "Custom": {
	//             "A": "value2",
	//             "B": 2
	//         }
	//     }
	// ]
}

func example2() {
	data := bytes.NewReader([]byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`))

	// Also supports `io.Reader`
	decoder := csv.NewDecoder(data)
	output := []Data{}
	err := decoder.Unmarshal(&data)
	if err != nil {
		panic(err)
	}
}
```

### Marshal

Marshal works in a very similar but the opposite way.
Takes in a slice or array of a struct, writes the header row of all the fields
then proceeds to write the contents of the slice as CSV data.

## TODO:

More tests need to be written, testing the usage of the built in CSV library's decoding/encoding capabilities.

Support more of the stdlib's types for marshalling and unmarshalling.
