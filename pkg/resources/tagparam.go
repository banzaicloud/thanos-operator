package resources

import (
	"fmt"
	"reflect"
	"strings"
)

func getArgs(input interface{}) []string {
	var args []string
	value := strctVal(input)
	fields := StructFields(value)
	for _, field := range fields {
		val := value.FieldByName(field.Name)
		tagArgFormat, _ := parseTagWithName(field.Tag.Get("thanos"))
		if tagArgFormat != "" {
			args = append(args, fmt.Sprintf(tagArgFormat, val))
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

// parseTag returns all the options in the tag
func parseTag(tag string) tagOptions {
	return tagOptions(strings.Split(tag, ","))
}

func StructFields(value reflect.Value) []reflect.StructField {
	t := value.Type()

	var f []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get("thanos"); tag == "-" {
			continue
		}

		f = append(f, field)
	}

	return f
}
