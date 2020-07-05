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

// UnmarshalCSV describes how CSV should handle types that aren't strings or built in
type UnmarshalCSV interface {
	UnmarshalCSV(string) error
}

var unmarshalCSV = reflect.TypeOf((*UnmarshalCSV)(nil)).Elem()

// Decoder reads and decodes a csv into an array from an input stream
type Decoder struct {
	reader *rawcsv.Reader
}

// NewDecoder creates a new decoder from the given reader
func NewDecoder(r io.Reader) *Decoder {
	reader := rawcsv.NewReader(r)
	return &Decoder{reader}
}

// SetDelimiter character for the csv reader
func (d *Decoder) SetDelimiter(delim rune) {
	d.reader.Comma = delim
}

// SetComment character for the csv reader
func (d *Decoder) SetComment(c rune) {
	d.reader.Comment = c
}

// TrimLeadingSpace the value for the underlying reader to true.
// If TrimLeadingSpace is true, leading white space in a field is ignored.
// This is done even if the field delimiter is white space.
func (d *Decoder) TrimLeadingSpace() {
	d.reader.TrimLeadingSpace = true
}

// DisableTrimLeadingSpace sets the value for the underlying reader to false.
// If TrimLeadingSpace is true, leading white space in a field is ignored.
// This is done even if the field delimiter is white space.
func (d *Decoder) DisableTrimLeadingSpace() {
	d.reader.TrimLeadingSpace = false
}

// Decode decodes the reader into the value v
// v must be an array of structs, where the struct field names (or tags) define the csv header name to decode from
// Will decode most built in types, otherwise it will use the FromString interface to decode
func (d *Decoder) Decode(v interface{}) error {
	value := reflect.ValueOf(v)

	if value.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Decode: could not decode into type %v - must be a pointer", value.Type())
	}

	ty := value.Type().Elem()

	switch ty.Kind() {
	case reflect.Array:
		// [n]struct { ... fields FromString }
		return fmt.Errorf("Decode: todo array types %v", ty)
	case reflect.Slice:
		// []struct { ... fields FromString }
		elem := ty.Elem()
		if elem.Kind() != reflect.Struct {
			return fmt.Errorf("Decode: could not decode into type %v - expected a slice of structs", ty)
		}

		fields := make([]string, 0, elem.NumField())
		for i := 0; i < cap(fields); i++ {
			field := elem.Field(i)

			column := field.Tag.Get("csv")
			if column == "" {
				column = field.Name
			}

			if !validUnmarshalType(field.Type) {
				return fmt.Errorf("Decode: %v is not a valid field type - try implement UnmarshalCSV for it", field.Type)
			}

			fields = append(fields, column)
		}

		headers, err := d.reader.Read()
		if err != nil {
			return err
		}

		h2f := map[int]int{} // headers to fields
		for i, header := range headers {
			for j, field := range fields {
				if header == field {
					h2f[i] = j
					goto next_header
				}
			}

			return fmt.Errorf("Decode: field for header[%s] was not found", header)

		next_header:
		}

		for {
			row, err := d.reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			record := reflect.New(elem)
			for i, column := range row {
				field := record.Elem().Field(h2f[i])

				if err := setField(field, column); err != nil {
					return err
				}
			}

			value.Elem().Set(reflect.Append(value.Elem(), record.Elem()))
		}

	case reflect.Map:
		// map[FromString][]FromString
		return fmt.Errorf("Decode: todo map types %v", ty)
	default:
		return fmt.Errorf("Decode: could not decode into type %v", ty)
	}

	return nil
}

func validUnmarshalType(ty reflect.Type) bool {
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

	if reflect.PtrTo(ty).Implements(unmarshalCSV) {
		return true
	}

	return false
}

func setField(field reflect.Value, value string) error {
	switch field.Type().Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
		return nil
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
		return nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(i)
		return nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(i)
		return nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(i)
		return nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
		return nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		field.SetFloat(f)
		return nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
		return nil
	case reflect.String:
		field.SetString(value)
		return nil
	}

	if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(t))
		return nil
	}

	return field.Addr().Interface().(UnmarshalCSV).UnmarshalCSV(value)
}

// Unmarshal the byte slice as a csv into the value v
// v must be an array of structs, where the struct field names (or tags) define the csv header name to decode from
// Will decode most built in types automatically, otherwise it will use the FromString interface to decode
func Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	d := NewDecoder(r)
	return d.Decode(v)
}
