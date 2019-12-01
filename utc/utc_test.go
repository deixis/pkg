package utc_test

import (
	"testing"
	"time"

	"github.com/deixis/pkg/utc"
)

// TestConversion ensures that a time.Time converted to UTC
// and back again to time.Time are the same
func TestConversion(t *testing.T) {
	tests := []time.Time{
		time.Now().UTC(),
		parse(t, "2006-01-02T15:04:05Z"),
		parse(t, "2017-03-09T00:00:00Z"),
		parse(t, "2053-03-09T00:00:00Z"),
		parse(t, "1970-01-01T00:00:00Z"),
		parse(t, "1900-12-31T00:00:00Z"),
	}

	for i, expect := range tests {
		u := utc.Convert(expect)
		got := u.Time()
		if expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, expect, got)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		In     []utc.UTC
		Expect utc.UTC
	}{
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "2006-01-02T15:04:05Z")),
				utc.Convert(parse(t, "2017-03-09T00:00:00Z")),
				utc.Convert(parse(t, "2053-03-09T00:00:00Z")),
			},
			Expect: utc.Convert(parse(t, "2006-01-02T15:04:05Z")),
		},
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "1900-01-01T00:00:00Z")),
				utc.Convert(parse(t, "1970-01-01T00:00:00Z")),
				utc.Convert(parse(t, "2040-01-01T00:00:00Z")),
			},
			Expect: utc.Convert(parse(t, "1900-01-01T00:00:00Z")),
		},
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "1900-01-01T00:00:03Z")),
				utc.Convert(parse(t, "1900-01-01T00:00:02Z")),
				utc.Convert(parse(t, "1900-01-01T00:00:01Z")),
			},
			Expect: utc.Convert(parse(t, "1900-01-01T00:00:01Z")),
		},
	}

	for i, test := range tests {
		got := utc.Min(test.In...)
		if test.Expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.Expect, got)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		In     []utc.UTC
		Expect utc.UTC
	}{
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "2006-01-02T15:04:05Z")),
				utc.Convert(parse(t, "2017-03-09T00:00:00Z")),
				utc.Convert(parse(t, "2053-03-09T00:00:00Z")),
			},
			Expect: utc.Convert(parse(t, "2053-03-09T00:00:00Z")),
		},
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "1900-01-01T00:00:00Z")),
				utc.Convert(parse(t, "1970-01-01T00:00:00Z")),
				utc.Convert(parse(t, "2040-01-01T00:00:00Z")),
			},
			Expect: utc.Convert(parse(t, "2040-01-01T00:00:00Z")),
		},
		{
			In: []utc.UTC{
				utc.Convert(parse(t, "1900-01-01T00:00:03Z")),
				utc.Convert(parse(t, "1900-01-01T00:00:02Z")),
				utc.Convert(parse(t, "1900-01-01T00:00:01Z")),
			},
			Expect: utc.Convert(parse(t, "1900-01-01T00:00:03Z")),
		},
	}

	for i, test := range tests {
		got := utc.Max(test.In...)
		if test.Expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.Expect, got)
		}
	}
}

// TestParse ensures that an RFC3339 string representation of time can be parsed
func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		expect time.Time
	}{
		{input: "2006-01-02T15:04:05Z", expect: parse(t, "2006-01-02T15:04:05Z")},
		{input: "2017-03-09T00:00:00Z", expect: parse(t, "2017-03-09T00:00:00Z")},
		{input: "2053-03-09T00:00:00Z", expect: parse(t, "2053-03-09T00:00:00Z")},
		{input: "1970-01-01T00:00:00Z", expect: parse(t, "1970-01-01T00:00:00Z")},
		{input: "1900-12-31T00:00:00Z", expect: parse(t, "1900-12-31T00:00:00Z")},
		{input: "2018-01-17T00:00:00+01:00", expect: parse(t, "2018-01-16T23:00:00Z")},
	}

	for i, test := range tests {
		got, err := utc.Parse(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if test.expect != got.Time() {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, got)
		}
	}
}

// TestIsZero ensures that the IsZero function behaves correctly
func TestIsZero(t *testing.T) {
	if utc.Now().IsZero() {
		t.Error("expect IsZero to not be true")
	}
}

// TestGobEncode ensures that time remains the same after being encoded/decoded in GOB
func TestGobEncode(t *testing.T) {
	tests := []time.Time{
		time.Now().UTC(),
		parse(t, "2006-01-02T15:04:05Z"),
		parse(t, "2017-03-09T00:00:00Z"),
		parse(t, "2053-03-09T00:00:00Z"),
		parse(t, "1970-01-01T00:00:00Z"),
		parse(t, "1900-12-31T00:00:00Z"),
	}

	for i, test := range tests {
		expect := utc.Convert(test)
		gob, err := expect.GobEncode()
		if err != nil {
			t.Fatal(err)
		}

		var got utc.UTC
		if err := got.GobDecode(gob); err != nil {
			t.Fatal(err)
		}

		if expect != got {
			t.Errorf("#%d - expect %d, but got %d", i, expect, got)
		}
	}
}

// TestMarshalJSON ensures that time remains the same after being marshalled/unmarshalled in JSON
func TestMarshalJSON(t *testing.T) {
	tests := []time.Time{
		time.Now().UTC(),
		parse(t, "2006-01-02T15:04:05Z"),
		parse(t, "2017-03-09T00:00:00Z"),
		parse(t, "2053-03-09T00:00:00Z"),
		parse(t, "1970-01-01T00:00:00Z"),
		parse(t, "1900-12-31T00:00:00Z"),
	}

	for i, test := range tests {
		expect := utc.Convert(test)
		json, err := expect.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		var got utc.UTC
		if err := got.UnmarshalJSON(json); err != nil {
			t.Fatal(err)
		}

		if expect != got {
			t.Errorf("#%d - expect %v, but got %v", i, expect.Time(), got.Time())
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		input    utc.UTC
		duration time.Duration
		expect   utc.UTC
	}{
		{
			input:    parseUTC(t, "2006-01-02T15:04:05Z"),
			duration: time.Hour,
			expect:   parseUTC(t, "2006-01-02T16:04:05Z"),
		},
		{
			input:    parseUTC(t, "2017-03-09T23:59:59Z"),
			duration: time.Second,
			expect:   parseUTC(t, "2017-03-10T00:00:00Z"),
		},
		{
			input:    parseUTC(t, "2053-03-09T00:00:00Z"),
			duration: time.Minute,
			expect:   parseUTC(t, "2053-03-09T00:01:00Z"),
		},
	}

	for i, test := range tests {
		got := test.input.Add(test.duration)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		input    utc.UTC
		duration time.Duration
		expect   utc.UTC
	}{
		{
			input:    parseUTC(t, "2006-01-02T15:04:05Z"),
			duration: time.Hour,
			expect:   parseUTC(t, "2006-01-02T14:04:05Z"),
		},
		{
			input:    parseUTC(t, "2017-03-10T00:00:00Z"),
			duration: time.Second,
			expect:   parseUTC(t, "2017-03-09T23:59:59Z"),
		},
		{
			input:    parseUTC(t, "2053-03-09T00:01:00Z"),
			duration: time.Minute,
			expect:   parseUTC(t, "2053-03-09T00:00:00Z"),
		},
	}

	for i, test := range tests {
		got := test.input.Sub(test.duration)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		a      utc.UTC
		b      utc.UTC
		expect time.Duration
	}{
		{
			a:      parseUTC(t, "2006-01-02T15:04:05Z"),
			b:      parseUTC(t, "2006-01-02T14:04:05Z"),
			expect: time.Hour,
		},
		{
			a:      parseUTC(t, "2017-03-10T00:00:00Z"),
			b:      parseUTC(t, "2017-03-09T23:59:59Z"),
			expect: time.Second,
		},
		{
			a:      parseUTC(t, "2053-03-09T00:01:00Z"),
			b:      parseUTC(t, "2053-03-09T00:00:00Z"),
			expect: time.Minute,
		},
	}

	for i, test := range tests {
		got := test.a.Distance(test.b)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, got)
		}
	}
}

func TestBeginningOfday(t *testing.T) {
	tests := []struct {
		input  utc.UTC
		expect utc.UTC
	}{
		{input: parseUTC(t, "2006-01-02T15:04:05Z"), expect: parseUTC(t, "2006-01-02T00:00:00Z")},
		{input: parseUTC(t, "2017-03-09T23:59:59Z"), expect: parseUTC(t, "2017-03-09T00:00:00Z")},
		{input: parseUTC(t, "2053-03-09T00:00:00Z"), expect: parseUTC(t, "2053-03-09T00:00:00Z")},
		{input: parseUTC(t, "1970-01-01T01:00:00Z"), expect: parseUTC(t, "1970-01-01T00:00:00Z")},
		{input: parseUTC(t, "1900-12-31T12:00:00Z"), expect: parseUTC(t, "1900-12-31T00:00:00Z")},
	}

	for i, test := range tests {
		got := test.input.BeginningOfDay()
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestEndOfday(t *testing.T) {
	tests := []struct {
		input  utc.UTC
		expect utc.UTC
	}{
		{input: parseUTC(t, "2006-01-02T15:04:05Z"), expect: parseUTC(t, "2006-01-02T23:59:59.999999999Z")},
		{input: parseUTC(t, "2017-03-09T23:59:59Z"), expect: parseUTC(t, "2017-03-09T23:59:59.999999999Z")},
		{input: parseUTC(t, "2053-03-09T00:00:00Z"), expect: parseUTC(t, "2053-03-09T23:59:59.999999999Z")},
		{input: parseUTC(t, "1970-01-01T01:00:00Z"), expect: parseUTC(t, "1970-01-01T23:59:59.999999999Z")},
		{input: parseUTC(t, "1900-12-31T12:00:00Z"), expect: parseUTC(t, "1900-12-31T23:59:59.999999999Z")},
	}

	for i, test := range tests {
		got := test.input.EndOfDay()
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		input  utc.UTC
		round  time.Duration
		expect utc.UTC
	}{
		{
			input:  parseUTC(t, "2006-01-02T15:04:05Z"),
			round:  time.Hour * 24,
			expect: parseUTC(t, "2006-01-02T00:00:00Z"),
		},
		{
			input:  parseUTC(t, "2017-03-09T23:59:59Z"),
			round:  time.Hour,
			expect: parseUTC(t, "2017-03-09T23:00:00Z"),
		},
		{
			input:  parseUTC(t, "2017-03-09T23:59:59Z"),
			round:  time.Minute * 5,
			expect: parseUTC(t, "2017-03-09T23:55:00Z"),
		},
		{
			input:  parseUTC(t, "1970-01-01T23:59:59.999999999Z"),
			round:  time.Second,
			expect: parseUTC(t, "1970-01-01T23:59:59Z"),
		},
		{
			input:  parseUTC(t, "1970-12-31T23:59:59.999999999Z"),
			round:  time.Nanosecond * 100,
			expect: parseUTC(t, "1970-12-31T23:59:59.9999999Z"),
		},
	}

	for i, test := range tests {
		got := test.input.Floor(test.round)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		input  utc.UTC
		round  time.Duration
		expect utc.UTC
	}{
		{
			input:  parseUTC(t, "2006-01-02T15:04:05Z"),
			round:  time.Hour * 24,
			expect: parseUTC(t, "2006-01-03T00:00:00Z"),
		},
		{
			input:  parseUTC(t, "2017-03-09T23:59:59Z"),
			round:  time.Hour,
			expect: parseUTC(t, "2017-03-10T00:00:00Z"),
		},
		{
			input:  parseUTC(t, "2017-03-09T23:59:59Z"),
			round:  time.Minute * 5,
			expect: parseUTC(t, "2017-03-10T00:00:00Z"),
		},
		{
			input:  parseUTC(t, "1970-12-31T23:49:00.999999999Z"),
			round:  time.Minute * 5,
			expect: parseUTC(t, "1970-12-31T23:50:00Z"),
		},
		{
			input:  parseUTC(t, "1970-01-01T23:59:59.999999999Z"),
			round:  time.Second,
			expect: parseUTC(t, "1970-01-02T00:00:00Z"),
		},
	}

	for i, test := range tests {
		got := test.input.Ceil(test.round)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect.RFC3339Nano(), got.RFC3339Nano())
		}
	}
}

func TestRFC3339Nano(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		input  utc.UTC
		expect string
	}{
		{input: parseUTC(t, "2006-01-02T23:59:59.000000000Z"), expect: "2006-01-02T23:59:59Z"},
		{input: parseUTC(t, "2053-03-09T23:59:59.123456789Z"), expect: "2053-03-09T23:59:59.123456789Z"},
		{input: parseUTC(t, "1970-01-01T23:59:59.999999999Z"), expect: "1970-01-01T23:59:59.999999999Z"},
		{input: parseUTC(t, now.Format(time.RFC3339Nano)), expect: now.Format(time.RFC3339Nano)},
	}

	for i, test := range tests {
		got := test.input.RFC3339Nano()
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, got)
		}
	}
}

func TestRFC3339(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		input  utc.UTC
		expect string
	}{
		{input: parseUTC(t, "2006-01-02T23:59:59.000000000Z"), expect: "2006-01-02T23:59:59Z"},
		{input: parseUTC(t, "2053-03-09T23:59:59.123456789Z"), expect: "2053-03-09T23:59:59Z"},
		{input: parseUTC(t, "1970-01-01T23:59:59.999999999Z"), expect: "1970-01-01T23:59:59Z"},
		{input: parseUTC(t, now.Format(time.RFC3339Nano)), expect: now.Format(time.RFC3339)},
	}

	for i, test := range tests {
		got := test.input.RFC3339()
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, got)
		}
	}
}

func parse(t *testing.T, s string) time.Time {
	res, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return res.UTC()
}

func parseUTC(t *testing.T, s string) utc.UTC {
	return utc.Convert(parse(t, s))
}
