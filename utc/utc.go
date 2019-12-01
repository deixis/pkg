// Package utc is a lightweit time struct stripped of its timezone awareness
//
// Known limitations : dates before the year ~1678 cannot be represented with this time struct,
// but most of the time we can live with it.
package utc

import (
	"math"
	"strconv"
	"time"
)

const day = time.Hour * 24

// UTC is the number of nanoseconds elapsed since January 1, 1970 UTC.
// The result is undefined if the Unix time in nanoseconds cannot be represented by an int64.
// Note that this means the result of calling UnixNano on the zero Time is undefined.
type UTC int64

// Convert converts a Time struct into UTC
func Convert(t time.Time) UTC {
	return UTC(t.UTC().UnixNano())
}

// Now returns the current time.
func Now() UTC {
	return Convert(time.Now())
}

// In returns the current time + d
func In(d time.Duration) UTC {
	return Convert(time.Now().Add(d))
}

// Max returns the latest time
func Max(l ...UTC) UTC {
	if len(l) == 0 {
		return UTC(0)
	}
	max := UTC(math.MinInt64)
	for _, t := range l {
		if t > max {
			max = t
		}
	}
	return max
}

// Min returns the earliest time
func Min(l ...UTC) UTC {
	if len(l) == 0 {
		return UTC(0)
	}
	min := UTC(math.MaxInt64)
	for _, t := range l {
		if t < min {
			min = t
		}
	}
	return min
}

// MustParse is like Parse, but it panics when there is a parsing error.
// It simplifies safe initialisation of UTC values.
func MustParse(s string) UTC {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}

// Parse parses a formatted string as defined by the RFC3339 and returns the time value it represents.
//
// Note: If it returns a time.ParseError with "second out of range" when the second is equal to 60,
// it means it is a leap second. The go time package does not handle leap seconds.
// More info: https://github.com/golang/go/issues/8728
// However, chances are so slim that it does not worth handling that scenario (for now)
func Parse(s string) (UTC, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return 0, err
	}
	return Convert(t), nil
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t UTC) IsZero() bool {
	return t == 0
}

// Time converts the UTC type to the standard time.Time
func (t UTC) Time() time.Time {
	return time.Unix(0, int64(t)).UTC()
}

// String returns the default UTC string representation
func (t UTC) String() string {
	return t.RFC3339Nano()
}

// RFC3339 returns a string representation in RFC3339 format
func (t UTC) RFC3339() string {
	return t.Time().Format(time.RFC3339)
}

// RFC3339Nano returns a string representation in RFC3339 format with nanoseconds
func (t UTC) RFC3339Nano() string {
	return t.Time().Format(time.RFC3339Nano)
}

// Add returns the time t+d
func (t UTC) Add(d time.Duration) UTC {
	return t + UTC(d)
}

// Sub returns the time t-d
func (t UTC) Sub(d time.Duration) UTC {
	return t - UTC(d)
}

// Distance returns the duration t - u.
// If the result exceeds the maximum (or minimum) value that can be stored in a Duration,
// the maximum (or minimum) duration will be returned.
func (t UTC) Distance(u UTC) time.Duration {
	return time.Duration(t) - time.Duration(u)
}

// BeginningOfDay returns a new UTC with its time reset to midnight
func (t UTC) BeginningOfDay() UTC {
	return t.Floor(day)
}

// EndOfDay returns a new UTC with its time set to the last nanosecond of the day
func (t UTC) EndOfDay() UTC {
	return t.Floor(day) + UTC(day-time.Nanosecond)
}

// Floor rounds date down to the given precision
func (t UTC) Floor(prec time.Duration) UTC {
	return t - abs(t)%UTC(prec)
}

// Ceil rounds date up to the given precision
func (t UTC) Ceil(prec time.Duration) UTC {
	rem := abs(t) % UTC(prec)
	if rem > 0 {
		return t - rem + UTC(prec)
	}
	return t - rem
}

// GobEncode implements the gob.GobEncoder interface.
func (t UTC) GobEncode() ([]byte, error) {
	s := strconv.FormatInt(int64(t), 10)
	return []byte(s), nil
}

// GobDecode implements the gob.GobDecoder interface.
func (t *UTC) GobDecode(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	v := UTC(i)
	*t = v
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t UTC) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UTC) UnmarshalJSON(data []byte) error {
	a := t.Time()
	if err := a.UnmarshalJSON(data); err != nil {
		return err
	}
	*t = Convert(a)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface
// It uses the default string representation from String()
func (t UTC) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
// It uses the default layout - RFC3339 with nanoseconds
func (t *UTC) UnmarshalText(text []byte) error {
	u, err := Parse(string(text))
	if err != nil {
		return err
	}
	*t = u
	return nil
}

func abs(t UTC) UTC {
	if t < 0 {
		return t * -1
	}
	return t
}
