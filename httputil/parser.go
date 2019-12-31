package httputil

import (
	"context"
	"encoding"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/deixis/errors"
	"github.com/deixis/spine/log"
	"github.com/deixis/spine/net/http"
)

const queryStringTag = "qs"

// ParseReq parses the request and returns a standard error in case of failure
func ParseReq(ctx context.Context, r *http.Request, params interface{}) error {
	err := r.Parse(ctx, params)
	if err != nil {
		log.Warn(ctx, "http.parse.err", "Cannot parse request",
			log.Error(err),
		)

		return errors.WithBad(err)
	}

	return nil
}

// ParseQuery parses the values of v from the HTTP query
func ParseQuery(q url.Values, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("httputil: ParseQuery(non-pointer " + reflect.TypeOf(v).String() + ")")
	}

	rv = reflect.Indirect(rv)
	tv := rv.Type()
	switch rv.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := tv.Field(i)
			val := rv.Field(i)

			// If the current field has a query string tags
			tag, ok := field.Tag.Lookup(queryStringTag)
			if !ok {
				continue
			}
			tag, opts := parseTag(tag)

			// If the query string has the given tag name
			qVal := q.Get(tag)
			if qVal == "" {
				if opts.Contains("required") {
					return errors.Bad(&errors.FieldViolation{
						Field:       tag,
						Description: "Missing query string",
					})
				}
				continue
			}
			// walk down to get the first non-pointer
			ut, nptr := indirect(val)
			if ut != nil {
				if err := ut.UnmarshalText([]byte(qVal)); err != nil {
					return errors.WithBad(err)
				}
				continue
			}
			val = nptr

			// Parse primitives
			switch reflect.Indirect(val).Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(qVal, 10, 64)
				if err != nil {
					return errors.WithBad(err)
				}
				val.SetInt(i)
			case reflect.Uint, reflect.Uint32, reflect.Uint64:
				i, err := strconv.ParseUint(qVal, 10, 64)
				if err != nil {
					return errors.WithBad(err)
				}
				val.SetUint(i)
			case reflect.Float32, reflect.Float64:
				f, err := strconv.ParseFloat(qVal, 64)
				if err != nil {
					return errors.WithBad(err)
				}
				val.SetFloat(f)
			case reflect.Bool:
				b, err := strconv.ParseBool(qVal)
				if err != nil {
					return errors.WithBad(err)
				}
				val.SetBool(b)
			case reflect.String:
				val.SetString(qVal)
			}
		}
	default:
		return errors.New("httputil: ParseQuery(unsupported type " + reflect.TypeOf(v).String() + ")")
	}

	return nil
}

// indirect walks down v allocating pointers as needed, until it gets to a non-pointer.
func indirect(v reflect.Value) (encoding.TextUnmarshaler, reflect.Value) {
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		v = v.Elem()
	}
	return nil, v
}

// tagOptions is the string following a comma in a struct field's "qs"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
