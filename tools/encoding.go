package tools

import (
	"strings"
	"unicode/utf8"
)

const LANG_EN = "eng"
const LANG_FR = "fre"

// Try to guess the encoding, depending on the lang, and convert s to utf8
func TryToUtf8(s string, lang string) string {

	var b = []byte(s)

	switch lang {
	case LANG_EN: // in English, no need to bother for accentuated characters
		return s
	case LANG_FR: // in French, we assume it's either utf8 or iso8859-1
		if !utf8.Valid(b) || strings.Count(s, string([]byte{0xef, 0xbf, 0xbd})) > 0 {
			return iso88591ToUtf8(b)
		}
	}
	return s
}

// Conversion from iso8859-1 to utf8
func iso88591ToUtf8(iso88591 []byte) string {
	buf := make([]rune, len(iso88591))
	for i, b := range iso88591 {
		buf[i] = rune(b)
	}
	return string(buf)
}
