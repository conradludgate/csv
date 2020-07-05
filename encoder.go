package csv

import (
	"bytes"
	rawcsv "encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

// MarshalCSV describes how CSV should handle types that aren't strings or built in
type MarshalCSV interface {
	MarshalCSV() string
}

var marshalCSV = reflect.TypeOf((*MarshalCSV)(nil)).Elem()

// Encoder encodes and writes the contents of a slice into a csv file
type Encoder struct {
	writer *rawcsv.Writer
}

// NewEncoder creates a new encoder from the given writer
func NewEncoder(w io.Writer) *Encoder {
	writer := rawcsv.NewWriter(w)
	return &Encoder{writer}
}

// SetDelimiter character for the csv reader
func (e *Encoder) SetDelimiter(delim rune) {
	e.writer.Comma = delim
}

// UseCRLF enables CRLF line ending for the csv Writer
func (e *Encoder) UseCRLF() {
	e.writer.UseCRLF = true
}

// NoCRLF disables CRLF line ending for the csv Writer
func (e *Encoder) NoCRLF() {
	e.writer.UseCRLF = false
}

// Encode and write the value of v into a csv
func (e *Encoder) Encode(v interface{}) error {
	value := reflect.ValueOf(v)
	ty := value.Type()

	for ty.Kind() == reflect.Ptr {
		value = value.Elem()
		ty = ty.Elem()
	}

	switch ty.Kind() {
	case reflect.Array, reflect.Slice:
		elem := ty.Elem()
		if elem.Kind() != reflect.Struct {
			return fmt.Errorf("Encode: could not encode type %v - expected a collection of structs", ty)
		}

		fl := elem.NumField()
		fields := make([]string, 0, fl)
		for i := 0; i < fl; i++ {
			field := elem.Field(i)

			column := field.Tag.Get("csv")
			if column == "" {
				column = field.Name
			}

			if !validMarshalType(field.Type) {
				return fmt.Errorf("Encode: %v is not a valid field type - try implement MarshalCSV for it", field.Type)
			}

			fields = append(fields, column)
		}

		if err := e.writer.Write(fields); err != nil {
			return err
		}

		l := value.Len()
		for i := 0; i < l; i++ {
			row := make([]string, fl)
			for j := 0; j < fl; j++ {
				row[j] = getValue(value.Index(i).Field(j))
			}
			if err := e.writer.Write(row); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Encode: could not encode type %v", ty)
	}

	e.writer.Flush()
	return e.writer.Error()
}

// Marshal the value v into a csv
func Marshal(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(b)
	err := encoder.Encode(v)
	return b.Bytes(), err
}

func validMarshalType(ty reflect.Type) bool {
	switch ty.Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	}

	if ty.PkgPath() == "time" && ty.Name() == "Time" {
		return true
	}

	if reflect.PtrTo(ty).Implements(marshalCSV) {
		return true
	}

	return false
}

func getValue(field reflect.Value) string {
	switch field.Type().Kind() {
	case reflect.Bool:
		if field.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(field.Float(), 'f', 6, 32)
	case reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', 15, 64)
	case reflect.String:
		return field.String()
	}

	if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
		return field.Interface().(time.Time).Format(time.RFC3339)
	}

	return field.Interface().(MarshalCSV).MarshalCSV()
}
