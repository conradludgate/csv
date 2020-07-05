package main

import (
	"encoding/json"
	"fmt"
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

func main() {
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
}
