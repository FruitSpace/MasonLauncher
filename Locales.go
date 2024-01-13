package main

import "github.com/jeandeaual/go-locale"

func getLocale() string {
	lang, err := locale.GetLanguage()
	if err != nil {
		lang = "en"
	}
	if lang != "ru" && lang != "ua" && lang != "kz" {
		lang = "en"
	} else {
		lang = "ru"
	}
	return lang
}

var lang = getLocale()

func GetLoc(q string) string {
	return loc[q+"_"+lang]
}

var loc = map[string]string{}
