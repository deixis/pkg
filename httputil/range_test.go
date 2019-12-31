package httputil_test

import (
	"testing"

	"github.com/deixis/pkg/httputil"
)

// TestParseContentRange ensures that Content-Range header are properly parsed
func TestParseContentRange(t *testing.T) {
	t.Parallel()

	table := []struct {
		in    string
		start int64
		end   int64
		size  int64
		err   error
	}{
		{
			in:    "bytes 0-63/128",
			start: 0,
			end:   63,
			size:  128,
		},
		{
			in:    "bytes 64-128/128",
			start: 64,
			end:   128,
			size:  128,
		},
		{
			in:  "",
			err: httputil.ErrInvalidFormat,
		},
		{
			in:  "bytes",
			err: httputil.ErrInvalidFormat,
		},
		{
			in:  "64-128/128",
			err: httputil.ErrInvalidFormat,
		},
		{
			in:  "bytes 64-128",
			err: httputil.ErrInvalidFormat,
		},
		{
			in:  "bytes 128/128",
			err: httputil.ErrInvalidFormat,
		},
		{
			in:  "bytes 129-193/128",
			err: httputil.ErrInvalidRange,
		},
	}

	for i, test := range table {
		res, err := httputil.ParseContentRange(test.in)
		if err != test.err {
			t.Errorf("%d - expect error <%s>, but got <%s>", i, test.err, err)
		}
		if test.err != nil {
			continue
		}

		if res.Start != test.start {
			t.Errorf("expect to get start %d, but got %d", test.start, res.Start)
		}
		if res.End != test.end {
			t.Errorf("expect to get end %d, but got %d", test.end, res.End)
		}
		if res.Size != test.size {
			t.Errorf("expect to get size %d, but got %d", test.size, res.Size)
		}
	}
}
