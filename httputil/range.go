package httputil

import (
	"errors"
	"strconv"
	"strings"
)

type HTTPRange struct {
	Start, End, Size int64
}

// ErrInvalidFormat is returned when it is an invalid range format
var ErrInvalidFormat = errors.New("invalid range format")

// ErrInvalidRange is returned when it is an invalid range
var ErrInvalidRange = errors.New("invalid range")

// ParseContentRange parses the HTTP Content-Range header
//
// e.g.
// Content-Range: <unit> <range-start>-<range-end>/<size>
// Content-Range: <unit> <range-start>-<range-end>/*
// Content-Range: <unit> */<size>
func ParseContentRange(s string) (*HTTPRange, error) {
	if s == "" {
		return nil, ErrInvalidFormat // header not present
	}

	// Only bytes are supported (for now)
	const b = "bytes "
	if !strings.HasPrefix(s, b) {
		return nil, ErrInvalidFormat
	}

	r := strings.Split(s[len(b):], "/")
	if len(r) != 2 {
		return nil, ErrInvalidFormat
	}
	ran := r[0]
	size := strings.TrimSpace(r[1])

	// Parse size
	i, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	httpRange := &HTTPRange{
		Size: i,
	}

	// Parse ranges
	r = strings.Split(ran, "-")
	if len(r) != 2 {
		return nil, ErrInvalidFormat
	}
	start := strings.TrimSpace(r[0])
	end := strings.TrimSpace(r[1])

	i, err = strconv.ParseInt(start, 10, 64)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	httpRange.Start = i
	i, err = strconv.ParseInt(end, 10, 64)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	httpRange.End = i

	// Ensure range validity
	if httpRange.Start > httpRange.End {
		return nil, ErrInvalidRange
	}
	if httpRange.End > httpRange.Size {
		return nil, ErrInvalidRange
	}

	return httpRange, nil
}
