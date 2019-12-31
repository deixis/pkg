package httputil_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/deixis/pkg/httputil"
	"github.com/deixis/pkg/utc"
)

type dummyParseQuery struct {
	Min          utc.UTC  `qs:"min"`
	Max          *utc.UTC `qs:"max"`
	Limit        uint     `qs:"limit"`
	Continuation string   `qs:"continuation"`
	Q            *string  `qs:"q"`
	B            bool     `qs:"b"`
	F            float64  `qs:"f"`
}

func TestParseQuery(t *testing.T) {
	t.Parallel()

	min := parseUTC(t, "2027-12-20T14:00:00Z")
	a := parseUTC(t, "2027-12-20T14:00:00Z")
	max := &a
	q := "foo"
	f := float64(3.141592653589793)

	table := []struct {
		input  url.Values
		expect dummyParseQuery
		err    bool
	}{
		{input: parseQuery(t, ""), expect: dummyParseQuery{}},
		{input: parseQuery(t, "min=2027-12-20T14:00:00Z"), expect: dummyParseQuery{Min: min}},
		{input: parseQuery(t, "max=2027-12-20T14:00:00Z"), expect: dummyParseQuery{Max: max}},
		{input: parseQuery(t, "limit=15"), expect: dummyParseQuery{Limit: 15}},
		{input: parseQuery(t, "continuation=g2gCbQAAAAdya"), expect: dummyParseQuery{Continuation: "g2gCbQAAAAdya"}},
		{input: parseQuery(t, "q=foo"), expect: dummyParseQuery{Q: &q}},
		{input: parseQuery(t, "min=2027-12-20T14:00:00Z&max=2027-12-20T14:00:00Z&limit=15&continuation=g2gCbQAAAAdya&q=foo&b=true&f=3.141592653589793"),
			expect: dummyParseQuery{
				Min:          min,
				Max:          max,
				Limit:        15,
				Continuation: "g2gCbQAAAAdya",
				Q:            &q,
				B:            true,
				F:            f,
			},
		},
		{input: parseQuery(t, "min=20"), err: true},
	}

	for i, test := range table {
		res := dummyParseQuery{}
		err := httputil.ParseQuery(test.input, &res)
		if (err != nil) != test.err {
			t.Errorf("#%d - expect to get error", i)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(test.expect, res) {
			t.Errorf("#%d - expect to get %v, but got %v", i, test.expect, res)
		}
	}
}

type dummyParseQueryRequired struct {
	Min utc.UTC  `qs:"min,required"`
	Max *utc.UTC `qs:"max"`
}

func TestParseQueryRequired(t *testing.T) {
	t.Parallel()

	min := parseUTC(t, "2027-12-20T14:00:00Z")
	max := parseUTC(t, "2027-12-20T14:00:00Z")

	table := []struct {
		input  url.Values
		expect dummyParseQueryRequired
		err    bool
	}{
		{input: parseQuery(t, ""), err: true},
		{input: parseQuery(t, "min=2027-12-20T14:00:00Z"), expect: dummyParseQueryRequired{Min: min}},
		{input: parseQuery(t, "max=2027-12-20T14:00:00Z"), err: true},
		{input: parseQuery(t, "min=2027-12-20T14:00:00Z&max=2027-12-20T14:00:00Z"), expect: dummyParseQueryRequired{Max: &max, Min: min}},
	}

	for i, test := range table {
		res := dummyParseQueryRequired{}
		err := httputil.ParseQuery(test.input, &res)
		if (err != nil) != test.err {
			t.Errorf("#%d - expect to get error", i)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(test.expect, res) {
			t.Errorf("#%d - expect to get %v, but got %v", i, test.expect, res)
		}
	}
}

func parseQuery(t *testing.T, s string) url.Values {
	v, err := url.ParseQuery(s)
	if err != nil {
		t.Fatal("parseQuery ", err)
	}
	return v
}

func parseUTC(t *testing.T, s string) utc.UTC {
	v, err := utc.Parse(s)
	if err != nil {
		t.Fatal("parseQuery ", err)
	}
	return v
}
