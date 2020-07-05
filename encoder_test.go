package csv

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (c Custom) MarshalCSV() string {
	return fmt.Sprintf("%s|%d", c.A, c.B)
}

func TestEncodePass(t *testing.T) {
	time1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:05-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2020-07-03T16:39:44+01:00")

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
