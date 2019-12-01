// Package lang parses, validates, and format language tags.
// Language is formatted according to the RFC 5646 (IETF language tag).
//
// Each language tag is composed of one or more "subtags" separated by hyphens (-).
// Each subtag is composed of basic Latin letters or digits only.
// Subtags occur in the following order:
//
//   1. A single primary language subtag based on a two-letter language code from ISO 639-1
//   2. Up to three optional extended language subtags composed of three letters each, separated by hyphens;
//		(There is currently no extended language subtag registered in the Language Subtag Registry
// 		without an equivalent and preferred primary language subtag
//   3. ...
//
// Examples:
//   * en
//   * en-GB
//   * fr-CH
//   * de-CH
//   * de-DE
package lang

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"
)

var (
	Mul             = Tag{T: language.MustParse("mul")}
	English         = Tag{T: language.MustParse("en")}
	French          = Tag{T: language.MustParse("fr")}
	German          = Tag{T: language.MustParse("de")}
	Italian         = Tag{T: language.MustParse("it")}
	BritishEnglish  = Tag{T: language.MustParse("en-GB")}
	AmericanEnglish = Tag{T: language.MustParse("en-US")}
	Spanish         = Tag{T: language.MustParse("es")}
	SpainSpanish    = Tag{T: language.MustParse("es-ES")}
	SwissFrench     = Tag{T: language.MustParse("fr-CH")}
	SwissGerman     = Tag{T: language.MustParse("de-CH")}
)

var (
	// ErrUnknownTag indicates that the language tag is well-formed, but unknown
	ErrUnknownTag = errors.New("unknown language tag")
	// ErrInvalidTag indicates that the language tag cannot be parsed
	ErrInvalidTag = errors.New("invalid language tag")
)

// Tag is a language tag
type Tag struct {
	T language.Tag
}

// Must panics when there is an error
func Must(t *Tag, err error) *Tag {
	if err != nil {
		panic(err)
	}
	return t
}

// Parse parses a 2- or 3-letter ISO 639 code.
// It returns an ErrInvalidLang if s is a well-formed but unknown language identifier
// or another error if another error occurred.
func Parse(s string) (*Tag, error) {
	t, err := language.Parse(s)
	switch err.(type) {
	case language.ValueError:
		return nil, ErrUnknownTag
	case nil:
	default:
		return nil, ErrInvalidTag
	}

	t, err = language.All.Canonicalize(t)
	if err != nil {
		return nil, ErrInvalidTag
	}

	return &Tag{T: t}, nil
}

// ParseAcceptLanguage parses the contents of a Accept-Language header as
// defined in http://www.ietf.org/rfc/rfc2616.txt
func ParseAcceptLanguage(s string) ([]*Tag, error) {
	tags, _, err := language.ParseAcceptLanguage(s)
	switch err.(type) {
	case language.ValueError:
		return nil, ErrUnknownTag
	case nil:
	default:
		return nil, ErrInvalidTag
	}

	l := make([]*Tag, len(tags))
	for i, tag := range tags {
		l[i] = &Tag{T: tag}
	}
	return l, nil
}

// Base returns the base language of the language tag. If the base language is
// unspecified, an attempt will be made to infer it from the context.
func (t Tag) Base() string {
	b, _ := t.T.Base()
	return b.String()
}

func (t Tag) String() string {
	return t.T.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Tag) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		s := string(data[1 : len(data)-1])
		tag, err := Parse(s)
		if err != nil {
			return fmt.Errorf("Error parsing string '%s': %s", s, err)
		}
		*t = *tag
		return nil
	}

	return fmt.Errorf("Error decoding string '%s'", data)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Tag) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

func (t Tag) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Tag) UnmarshalText(p []byte) error {
	tag, err := Parse(string(p))
	if err != nil {
		return err
	}
	*t = *tag
	return nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (t Tag) GobEncode() ([]byte, error) {
	return []byte(t.String()), nil
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (t *Tag) GobDecode(data []byte) error {
	tag, err := Parse(string(data))
	if err != nil {
		return err
	}
	*t = *tag
	return nil
}

func (t Tag) DeepCopy(dst interface{}) error {
	switch dst := dst.(type) {
	case *Tag:
		dst.T = language.Make(t.T.String())
		return nil
	}
	return fmt.Errorf("Tag deep copy on an unknown type %T", dst)
}

type Matcher struct {
	supported []Tag
	matcher   language.Matcher
}

func NewMatcher(supported []Tag) *Matcher {
	l := make([]language.Tag, len(supported))
	for i, t := range supported {
		l[i] = t.T
	}
	return &Matcher{
		supported: supported,
		matcher:   language.NewMatcher(l),
	}
}

// Match finds the best supported language based on the preferred list and
// the languages for which there exists translations
func (m *Matcher) Match(preferred ...*Tag) *Tag {
	var l []language.Tag
	for _, t := range preferred {
		if t != nil {
			l = append(l, t.T)
		}
	}
	tag, i, _ := m.matcher.Match(l...)
	if i >= len(m.supported) {
		return &Tag{T: tag}
	}
	t := m.supported[i]
	return &t
}
