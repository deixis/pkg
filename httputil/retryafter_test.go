package httputil_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/deixis/pkg/httputil"
)

func TestParseRetryAfter(t *testing.T) {
	t.Parallel()

	httputil.Now = func() time.Time {
		return time.Date(2015, 10, 21, 7, 28, 15, 0, time.UTC)
	}

	table := []struct {
		input  http.Header
		expect time.Duration
		ok     bool
	}{
		{input: makeHeader("Retry-After", "10"), expect: 10 * time.Second, ok: true},
		{input: makeHeader("Retry-After", "Wed, 21 Oct 2015 07:28:00 GMT"), expect: 15 * time.Second, ok: true},
		{input: makeHeader("Retry-After", "0"), expect: 0, ok: true},
		{input: makeHeader("Retry-After", ""), expect: 0, ok: false},
		{input: makeHeader("Retry-After", "-10"), expect: 0, ok: true},
		{input: makeHeader("Retry-After", "ABC"), expect: 0, ok: false},
	}

	for i, test := range table {
		got, ok := httputil.ParseRetryAfter(test.input)
		if test.ok != ok {
			t.Errorf("#%d - expect ok %t, but got %t", i, test.ok, ok)
		}
		if test.expect != got {
			t.Errorf("#%d - expect duration %s, but got %s", i, test.expect, got)
		}
	}
}

func TestFormatRetryAfter(t *testing.T) {
	t.Parallel()

	httputil.Now = func() time.Time {
		return time.Date(2015, 10, 21, 7, 28, 15, 0, time.UTC)
	}

	table := []struct {
		input  time.Duration
		expect string
	}{
		{input: time.Minute, expect: "60"},
		{input: time.Hour, expect: "3600"},
		{input: 0, expect: "0"},
		{input: -60, expect: "0"},
	}

	for i, test := range table {
		h := http.Header{}
		httputil.FormatRetryAfter(h, test.input)

		got := h.Get("Retry-After")
		if test.expect != got {
			t.Errorf("#%d - expect `Retry-After` to be %s, but got %s", i, test.expect, got)
		}
	}
}

func makeHeader(k, v string) http.Header {
	h := http.Header{}
	h.Set(k, v)
	return h
}
