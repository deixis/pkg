package utc_test

import (
	"testing"

	"github.com/deixis/pkg/lang"
	"github.com/deixis/pkg/utc"
)

func TestFormats(t *testing.T) {
	tests := []struct {
		input  utc.UTC
		fmt    func(utc.UTC, lang.Tag) string
		lang   lang.Tag
		expect string
	}{
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateLong,
			lang:   lang.English,
			expect: "January 2, 2006",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateShort,
			lang:   lang.English,
			expect: "1/2/06",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateLong,
			lang:   lang.French,
			expect: "2 janvier 2006",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateShort,
			lang:   lang.French,
			expect: "02/01/2006",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateShort,
			lang:   lang.SwissFrench,
			expect: "02.01.06",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateLong,
			lang:   lang.German,
			expect: "2. Januar 2006",
		},
		{
			input:  parseUTC(t, "2006-01-02T00:00:00.000000000Z"),
			fmt:    utc.FormatDateShort,
			lang:   lang.German,
			expect: "02.01.06",
		},
		{
			input:  parseUTC(t, "2006-01-02T15:30:45.000000000Z"),
			fmt:    utc.FormatTimeShort,
			lang:   lang.German,
			expect: "15:30",
		},
	}

	for i, test := range tests {
		got := test.fmt(test.input, test.lang)
		if test.expect != got {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, got)
		}
	}
}
