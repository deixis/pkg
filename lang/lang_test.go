package lang_test

import (
	"reflect"
	"testing"

	"github.com/deixis/pkg/lang"
)

func TestParse(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect lang.Tag
		err    error
	}{
		{input: "en", expect: lang.English},
		{input: "en-GB", expect: lang.BritishEnglish},
		{input: "en_GB", expect: lang.BritishEnglish},
		{input: "fr-CH", expect: lang.SwissFrench},
		{input: "de-CH", expect: lang.SwissGerman},
		{input: "xx", err: lang.ErrUnknownTag},
		{input: "", err: lang.ErrInvalidTag},
	}

	for i, test := range table {
		res, err := lang.Parse(test.input)
		if err != nil {
			if err != test.err {
				t.Errorf("#%d expect error %s, but got %s", i, test.err, err)
			}
			continue
		}

		if test.expect != *res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

// TestString ensures that a parsed language tag has the same string representation
func TestString(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"en",
		"en-GB",
		"en-US",
		"fr-CH",
		"fr",
		"de-CH",
	}

	// String
	for i, input := range inputs {
		res, err := lang.Parse(input)
		if err != nil {
			t.Errorf("#%d unexpected error %s", i, err)
			continue
		}

		got := res.String()
		if got != input {
			t.Errorf("#%d - expect %s, but got %s", i, input, got)
		}
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect []*lang.Tag
		err    error
	}{
		{
			input:  "en-GB,en-US;q=0.8,en;q=0.6",
			expect: []*lang.Tag{&lang.BritishEnglish, &lang.AmericanEnglish, &lang.English},
		},
		{
			input:  "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			expect: []*lang.Tag{&lang.SwissFrench, &lang.French, &lang.English, &lang.German, &lang.Mul},
		},
		{
			input:  "",
			expect: []*lang.Tag{},
		},
		{
			input:  "es-ES, es;q=0.9, en-GB;q=0.8, en;q=0.7, de;q=0.6, *;q=0.5",
			expect: []*lang.Tag{&lang.SpainSpanish, &lang.Spanish, &lang.BritishEnglish, &lang.English, &lang.German, &lang.Mul},
		},
		{
			input:  "es-ES, es;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			expect: []*lang.Tag{&lang.SpainSpanish, &lang.Spanish, &lang.English, &lang.German, &lang.Mul},
		},
		{
			input:  "*;q=0.5",
			expect: []*lang.Tag{&lang.Mul},
		},
		{
			input:  "es-ES, es;q=0.9",
			expect: []*lang.Tag{&lang.SpainSpanish, &lang.Spanish},
		},
		{
			input:  "de-CH",
			expect: []*lang.Tag{&lang.SwissGerman},
		},
		{
			input: "xx",
			err:   lang.ErrUnknownTag,
		},
		{
			input: "obviously wrong",
			err:   lang.ErrInvalidTag,
		},
	}

	for i, test := range table {
		res, err := lang.ParseAcceptLanguage(test.input)
		if err != nil {
			if err != test.err {
				t.Errorf("#%d expect error %s, but got %s", i, test.err, err)
			}
			continue
		}

		if !reflect.DeepEqual(test.expect, res) {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

func TestMatcher_Match(t *testing.T) {
	t.Parallel()

	supported := []lang.Tag{
		lang.English,
		lang.BritishEnglish,
		lang.SwissGerman,
		lang.SwissFrench,
		lang.AmericanEnglish,
	}
	matcher := lang.NewMatcher(supported)

	table := []struct {
		input  *lang.Tag
		expect string
	}{
		{input: lang.Must(lang.Parse("en")), expect: "en"},
		{input: lang.Must(lang.Parse("fr")), expect: "fr-CH"},
		{input: lang.Must(lang.Parse("de")), expect: "de-CH"},
		{input: lang.Must(lang.Parse("en-GB")), expect: "en-GB"},
		{input: lang.Must(lang.Parse("fr-CH")), expect: "fr-CH"},
		{input: lang.Must(lang.Parse("de-CH")), expect: "de-CH"},
		{input: lang.Must(lang.Parse("it")), expect: "en"},
		{input: lang.Must(lang.Parse("it-CH")), expect: "en"},
		{input: lang.Must(lang.Parse("it-CH")), expect: "en"},
		{input: lang.Must(lang.Parse("es")), expect: "en"},
	}

	for i, test := range table {
		res := matcher.Match(test.input)
		if test.expect != res.String() {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res.String())
		}
	}

	res := lang.Must(lang.Parse("en"))
	if *res != lang.English {
		t.Error("expect to get fallback language")
	}
}

func TestMatcher_Empty(t *testing.T) {
	t.Parallel()

	supported := []lang.Tag{}
	matcher := lang.NewMatcher(supported)

	table := []struct {
		input  *lang.Tag
		expect string
	}{
		{input: lang.Must(lang.Parse("en")), expect: "und"},
		{input: lang.Must(lang.Parse("fr")), expect: "und"},
		{input: lang.Must(lang.Parse("de")), expect: "und"},
	}

	for i, test := range table {
		res := matcher.Match(test.input)
		if test.expect != res.String() {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res.String())
		}
	}

	res := lang.Must(lang.Parse("en"))
	if *res != lang.English {
		t.Error("expect to get fallback language")
	}
}

func TestTag_Base(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  lang.Tag
		expect lang.Tag
	}{
		{input: lang.BritishEnglish, expect: lang.English},
		{input: lang.English, expect: lang.English},
		{input: lang.SwissFrench, expect: lang.French},
		{input: lang.SwissGerman, expect: lang.German},
		{input: lang.German, expect: lang.German},
		{input: lang.Italian, expect: lang.Italian},
	}

	for i, test := range table {
		res := test.input.Base()
		if test.expect.String() != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

func TestTag_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	table := []struct {
		input lang.Tag
	}{
		{input: lang.BritishEnglish},
		{input: lang.English},
		{input: lang.SwissFrench},
		{input: lang.SwissGerman},
		{input: lang.German},
		{input: lang.Italian},
	}

	for i, test := range table {
		data, err := test.input.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		var res lang.Tag
		if err := res.UnmarshalJSON(data); err != nil {
			t.Fatal(err)
		}

		if test.input != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}

func TestTag_GobEncode(t *testing.T) {
	t.Parallel()

	table := []struct {
		input lang.Tag
	}{
		{input: lang.BritishEnglish},
		{input: lang.English},
		{input: lang.SwissFrench},
		{input: lang.SwissGerman},
		{input: lang.German},
		{input: lang.Italian},
	}

	for i, test := range table {
		data, err := test.input.GobEncode()
		if err != nil {
			t.Fatal(err)
		}

		var res lang.Tag
		if err := res.GobDecode(data); err != nil {
			t.Fatal(err)
		}

		if test.input != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}

func TestTag_MarshalText(t *testing.T) {
	t.Parallel()

	table := []struct {
		input lang.Tag
	}{
		{input: lang.BritishEnglish},
		{input: lang.English},
		{input: lang.SwissFrench},
		{input: lang.SwissGerman},
		{input: lang.German},
		{input: lang.Italian},
	}

	for i, test := range table {
		data, err := test.input.MarshalText()
		if err != nil {
			t.Fatal(err)
		}

		var res lang.Tag
		if err := res.UnmarshalText(data); err != nil {
			t.Fatal(err)
		}

		if test.input != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}

func TestTag_DeepCopy(t *testing.T) {
	t.Parallel()

	table := []struct {
		input lang.Tag
	}{
		{input: lang.BritishEnglish},
		{input: lang.English},
		{input: lang.SwissFrench},
		{input: lang.SwissGerman},
		{input: lang.German},
		{input: lang.Italian},
	}

	for i, test := range table {
		got := lang.Tag{}
		err := test.input.DeepCopy(&got)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(test.input, got) {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, got)
		}
	}
}
