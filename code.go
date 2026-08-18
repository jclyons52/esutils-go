package esutils

// Character predicates, mirroring esutils lib/code.js. `ch` is a code point.

// IsDecimalDigit reports whether ch is 0..9.
func IsDecimalDigit(ch int) bool { return 0x30 <= ch && ch <= 0x39 }

// IsHexDigit reports whether ch is a hex digit 0..9, a..f, A..F.
func IsHexDigit(ch int) bool {
	return (0x30 <= ch && ch <= 0x39) ||
		(0x61 <= ch && ch <= 0x66) ||
		(0x41 <= ch && ch <= 0x46)
}

// IsOctalDigit reports whether ch is 0..7.
func IsOctalDigit(ch int) bool { return 0x30 <= ch && ch <= 0x37 }

// nonASCIIWhiteSpaces are the white-space code points at/above U+1680
// (ECMA-262 7.2), matching esutils' NON_ASCII_WHITESPACES.
var nonASCIIWhiteSpaces = []int{
	0x1680,
	0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007,
	0x2008, 0x2009, 0x200A, 0x202F, 0x205F,
	0x3000,
	0xFEFF,
}

// IsWhiteSpace reports whether ch is ECMAScript white space (7.2).
func IsWhiteSpace(ch int) bool {
	if ch == 0x20 || ch == 0x09 || ch == 0x0B || ch == 0x0C || ch == 0xA0 {
		return true
	}
	if ch < 0x1680 {
		return false
	}
	for _, w := range nonASCIIWhiteSpaces {
		if ch == w {
			return true
		}
	}
	return false
}

// IsLineTerminator reports whether ch is an ECMAScript line terminator (7.3).
func IsLineTerminator(ch int) bool {
	return ch == 0x0A || ch == 0x0D || ch == 0x2028 || ch == 0x2029
}

// ASCII identifier tables (7.6), matching esutils' IDENTIFIER_START/PART.
var asciiIdentStart, asciiIdentPart [0x80]bool

func init() {
	for ch := 0; ch < 0x80; ch++ {
		asciiIdentStart[ch] = (0x61 <= ch && ch <= 0x7A) || // a..z
			(0x41 <= ch && ch <= 0x5A) || // A..Z
			ch == 0x24 || ch == 0x5F // $ _
		asciiIdentPart[ch] = asciiIdentStart[ch] || (0x30 <= ch && ch <= 0x39)
	}
}

// IsIdentifierStartES5 reports whether ch may start an ES5 identifier.
func IsIdentifierStartES5(ch int) bool {
	if ch < 0x80 {
		return asciiIdentStart[ch]
	}
	return ES5Start.contains(ch)
}

// IsIdentifierPartES5 reports whether ch may continue an ES5 identifier.
func IsIdentifierPartES5(ch int) bool {
	if ch < 0x80 {
		return asciiIdentPart[ch]
	}
	return ES5Part.contains(ch)
}

// IsIdentifierStartES6 reports whether ch may start an ES6 identifier.
func IsIdentifierStartES6(ch int) bool {
	if ch < 0x80 {
		return asciiIdentStart[ch]
	}
	return ES6Start.contains(ch)
}

// IsIdentifierPartES6 reports whether ch may continue an ES6 identifier.
func IsIdentifierPartES6(ch int) bool {
	if ch < 0x80 {
		return asciiIdentPart[ch]
	}
	return ES6Part.contains(ch)
}
