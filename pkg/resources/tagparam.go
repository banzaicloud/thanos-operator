package resources

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetArgs(input interface{}) []string {
	var args []string
	value := strctVal(input)
	elements := StructElements(value)
	for _, element := range elements {
		var val string
		tagArgFormat, _ := parseTagWithName(element.Field.Tag.Get("thanos"))
		if tagArgFormat != "" {
			switch i := element.Value.Interface().(type) {
			case v1.Duration:
				if i.Duration != 0 {
					val = i.String()
				}
			case int:
				if i != 0 {
					strconv.Itoa(i)
				}
			case string:
				val = i
			case bool:
				// Bool params are switches don't need to render value
				if i {
					args = append(args, tagArgFormat)
				}
			default:
				val = ""
			}
			if val != "" {
				args = append(args, fmt.Sprintf(tagArgFormat, val))
			}
		}
	}
	return args
}

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func parseTagWithName(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

// tagOptions contains a slice of tag options
type tagOptions []string

// Has returns true if the given option is available in tagOptions
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// Has returns true if the given option is available in tagOptions
func (t tagOptions) ValueForPrefix(opt string) (bool, string) {
	for _, tagOpt := range t {
		if strings.HasPrefix(tagOpt, opt) {
			return true, strings.Replace(tagOpt, opt, "", 1)
		}
	}
	return false, ""
}

type StructElement struct {
	Field reflect.StructField
	Value reflect.Value
}

func StructElements(value reflect.Value) []StructElement {
	t := value.Type()

	var f []StructElement

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get("thanos"); tag == "-" {
			continue
		}
		f = append(f, StructElement{
			Field: field,
			Value: fieldValue,
		})
	}

	return f
}
