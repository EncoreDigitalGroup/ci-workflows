package main

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// String provides static string manipulation utilities
var String = struct {
	TitleCase   func(string) string
	Language    func() language.Tag
	SetLanguage func(language.Tag)
}{
	TitleCase:   titleCase,
	Language:    getLanguage,
	SetLanguage: setLanguage,
}

// Private package variables
var (
	defaultLanguage = language.English
	lang            = defaultLanguage
	titleCaser      = cases.Title(defaultLanguage)
)

// Private implementation functions
func titleCase(str string) string {
	return titleCaser.String(str)
}

func getLanguage() language.Tag {
	return lang
}

func setLanguage(newLang language.Tag) {
	lang = newLang
	titleCaser = cases.Title(newLang)
}
