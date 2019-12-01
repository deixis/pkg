package utc

import (
	"github.com/deixis/pkg/lang"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/de_CH"
	"github.com/go-playground/locales/de_DE"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/en_GB"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/fr_CH"
	"github.com/go-playground/locales/fr_FR"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/it_CH"
	"github.com/go-playground/locales/it_IT"
)

var localeMapper = map[string]locales.Translator{
	"de":    de.New(),
	"de-CH": de_CH.New(),
	"de-DE": de_DE.New(),
	"en":    en.New(),
	"en-GB": en_GB.New(),
	"en-US": en_US.New(),
	"fr":    fr.New(),
	"fr-CH": fr_CH.New(),
	"fr-FR": fr_FR.New(),
	"it":    it.New(),
	"it-CH": it_IT.New(),
	"it-IT": it_CH.New(),
}

var fallback = en.New()

func FormatDateLong(u UTC, lang lang.Tag) string {
	if l, ok := localeMapper[lang.String()]; ok {
		return l.FmtDateLong(u.Time())
	}
	return fallback.FmtDateLong(u.Time())
}

func FormatDateShort(u UTC, lang lang.Tag) string {
	if l, ok := localeMapper[lang.String()]; ok {
		return l.FmtDateShort(u.Time())
	}
	return fallback.FmtDateShort(u.Time())
}

func FormatTimeShort(u UTC, lang lang.Tag) string {
	if l, ok := localeMapper[lang.String()]; ok {
		return l.FmtTimeShort(u.Time())
	}
	return fallback.FmtTimeShort(u.Time())
}
