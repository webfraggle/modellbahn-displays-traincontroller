package api

import (
	"strings"
	"unicode/utf8"
)

// repairUTF8 fixes double-encoded UTF-8 characters that occur when
// TrainController (Windows) passes strings that were mis-interpreted as
// Windows-1252 and then re-encoded as UTF-8.
func repairUTF8(s string) string {
	for broken, fixed := range replacements {
		s = strings.ReplaceAll(s, broken, fixed)
	}
	return s
}

// replacements is built at startup by computeReplacements.
var replacements map[string]string

func init() {
	replacements = computeReplacements()
}

// win1252Special maps Windows-1252 bytes 0x80–0x9F to their Unicode codepoints.
// A zero entry means the byte is undefined in Windows-1252.
var win1252Special = [32]rune{
	0x20AC, 0x0000, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021,
	0x02C6, 0x2030, 0x0160, 0x2039, 0x0152, 0x0000, 0x017D, 0x0000,
	0x0000, 0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014,
	0x02DC, 0x2122, 0x0161, 0x203A, 0x0153, 0x0000, 0x017E, 0x0178,
}

// computeReplacements builds a map from "broken" UTF-8 strings to their
// correct Unicode equivalents.
//
// A "broken" string occurs when a Unicode character's UTF-8 byte sequence is
// read byte-by-byte and each byte is converted to its Windows-1252 character,
// then those characters are UTF-8 encoded together. The result looks like
// garbled text (e.g. "Ã¤" instead of "ä").
//
// We compute the broken form of every character in U+00A0–U+00FF (Latin-1
// supplement) and every character in the Windows-1252 special range (U+20AC
// and similar), then map each broken form back to the correct character.
func computeReplacements() map[string]string {
	m := make(map[string]string)

	addIfMultibyte := func(cp rune) {
		broken := brokenUTF8(cp)
		if broken != "" {
			m[broken] = string(cp)
		}
	}

	// Latin-1 supplement: U+00A0–U+00FF
	for cp := rune(0x00A0); cp <= 0x00FF; cp++ {
		addIfMultibyte(cp)
	}

	// Windows-1252 special characters (U+0080–U+009F range in Windows-1252)
	for _, cp := range win1252Special {
		if cp != 0 {
			addIfMultibyte(cp)
		}
	}

	return m
}

// brokenUTF8 returns the garbled string that results from taking the UTF-8
// encoding of cp and misinterpreting each byte as a Windows-1252 character.
// Returns "" if cp is a single-byte character (no garbling possible) or if
// any byte maps to an undefined Windows-1252 codepoint.
func brokenUTF8(cp rune) string {
	var buf [4]byte
	n := utf8.EncodeRune(buf[:], cp)
	if n <= 1 {
		return ""
	}

	result := make([]rune, 0, n)
	for _, b := range buf[:n] {
		var r rune
		if b >= 0x80 && b <= 0x9F {
			r = win1252Special[b-0x80]
			if r == 0 {
				// Undefined Windows-1252 byte: treat as the bare Unicode codepoint
				// (U+0080–U+009F), which is how many implementations handle it.
				r = rune(b)
			}
		} else {
			r = rune(b) // 0x00–0x7F and 0xA0–0xFF are identical in Latin-1 and Unicode
		}
		result = append(result, r)
	}

	return string(result)
}
