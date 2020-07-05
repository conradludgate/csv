package csv_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/conradludgate/csv"
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

func (c Custom) MarshalCSV() string {
	return fmt.Sprintf("%s|%d", c.A, c.B)
}

func ExampleUnmarshal() {
	data := []byte(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	output := []Data{}
	err := csv.Unmarshal(data, &output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleDecoder_Decode() {
	data := strings.NewReader(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	decoder := csv.NewDecoder(data)
	output := []Data{}
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleDecoder_SetDelimiter() {
	data := strings.NewReader(`Foo~bar~Time~Custom
hello world~9223372036854775807~2006-01-02T15:04:05-07:00~value1|1
goodbye world~-9223372036854775808~2020-07-03T16:39:44+01:00~value2|2`)

	decoder := csv.NewDecoder(data)
	decoder.SetDelimiter('~')
	output := []Data{}
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleDecoder_SetComment() {
	data := strings.NewReader(`Foo,bar,Time,Custom
hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
# This line is a comment and will be ignored
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	decoder := csv.NewDecoder(data)
	decoder.SetComment('#')
	output := []Data{}
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleDecoder_TrimLeadingSpace() {
	data := strings.NewReader(`          Foo,                 bar,                     Time,  Custom
  hello world, 9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	decoder := csv.NewDecoder(data)
	decoder.TrimLeadingSpace()
	output := []Data{}
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleDecoder_DisableTrimLeadingSpace() {
	data := strings.NewReader(`Foo,bar,Time,Custom
  hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2`)

	decoder := csv.NewDecoder(data)
	decoder.DisableTrimLeadingSpace()
	output := []Data{}
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}

	outputJSON, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(outputJSON))
	// Output:
	// [
	// 	{
	// 		"Foo": "  hello world",
	// 		"Bar": 9223372036854775807,
	// 		"Time": "2006-01-02T15:04:05-07:00",
	// 		"Custom": {
	// 			"A": "value1",
	// 			"B": 1
	// 		}
	// 	},
	// 	{
	// 		"Foo": "goodbye world",
	// 		"Bar": -9223372036854775808,
	// 		"Time": "2020-07-03T16:39:44+01:00",
	// 		"Custom": {
	// 			"A": "value2",
	// 			"B": 2
	// 		}
	// 	}
	// ]
}

func ExampleMarshal() {
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

	bytes, _ := csv.Marshal(data)
	fmt.Println(string(bytes))
	// Output:
	// Foo,bar,Time,Custom
	// hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
	// goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2
}

func ExampleEncoder_Encode() {
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

	encoder := csv.NewEncoder(os.Stdout)
	encoder.Encode(data)
	// Output:
	// Foo,bar,Time,Custom
	// hello world,9223372036854775807,2006-01-02T15:04:05-07:00,value1|1
	// goodbye world,-9223372036854775808,2020-07-03T16:39:44+01:00,value2|2
}

func ExampleEncoder_SetDelimiter() {
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

	encoder := csv.NewEncoder(os.Stdout)
	encoder.SetDelimiter('~')
	encoder.Encode(data)
	// Output:
	// Foo~bar~Time~Custom
	// hello world~9223372036854775807~2006-01-02T15:04:05-07:00~value1|1
	// goodbye world~-9223372036854775808~2020-07-03T16:39:44+01:00~value2|2
}
