package esutils

// Keyword and identifier-name predicates, mirroring esutils lib/keyword.js.

// isStrictModeReservedWordES6 mirrors the ES6 strict-mode reserved words.
func isStrictModeReservedWordES6(id string) bool {
	switch id {
	case "implements", "interface", "package", "private", "protected",
		"public", "static", "let":
		return true
	}
	return false
}

// IsKeywordES5 reports whether id is an ES5 keyword (yield exempt in sloppy
// mode), matching esutils keyword.js.
func IsKeywordES5(id string, strict bool) bool {
	if !strict && id == "yield" {
		return false
	}
	return IsKeywordES6(id, strict)
}

// IsKeywordES6 reports whether id is an ES6 keyword.
func IsKeywordES6(id string, strict bool) bool {
	if strict && isStrictModeReservedWordES6(id) {
		return true
	}
	switch len(id) {
	case 2:
		return id == "if" || id == "in" || id == "do"
	case 3:
		return id == "var" || id == "for" || id == "new" || id == "try"
	case 4:
		return id == "this" || id == "else" || id == "case" ||
			id == "void" || id == "with" || id == "enum"
	case 5:
		return id == "while" || id == "break" || id == "catch" ||
			id == "throw" || id == "const" || id == "yield" ||
			id == "class" || id == "super"
	case 6:
		return id == "return" || id == "typeof" || id == "delete" ||
			id == "switch" || id == "export" || id == "import"
	case 7:
		return id == "default" || id == "finally" || id == "extends"
	case 8:
		return id == "function" || id == "continue" || id == "debugger"
	case 10:
		return id == "instanceof"
	}
	return false
}

// IsReservedWordES5 reports whether id is an ES5 reserved word.
func IsReservedWordES5(id string, strict bool) bool {
	return id == "null" || id == "true" || id == "false" || IsKeywordES5(id, strict)
}

// IsReservedWordES6 reports whether id is an ES6 reserved word.
func IsReservedWordES6(id string, strict bool) bool {
	return id == "null" || id == "true" || id == "false" || IsKeywordES6(id, strict)
}

// IsRestrictedWord reports whether id is eval or arguments.
func IsRestrictedWord(id string) bool {
	return id == "eval" || id == "arguments"
}

// IsIdentifierNameES5 reports whether id is a valid ES5 identifier name.
func IsIdentifierNameES5(id string) bool {
	if len(id) == 0 {
		return false
	}
	runes := []rune(id)
	if !IsIdentifierStartES5(int(runes[0])) {
		return false
	}
	for _, r := range runes[1:] {
		if !IsIdentifierPartES5(int(r)) {
			return false
		}
	}
	return true
}

// IsIdentifierNameES6 reports whether id is a valid ES6 identifier name
// (supplementary-plane scalar values decoded from surrogate pairs).
func IsIdentifierNameES6(id string) bool {
	if len(id) == 0 {
		return false
	}
	for i, r := range []rune(id) {
		if i == 0 {
			if !IsIdentifierStartES6(int(r)) {
				return false
			}
		} else if !IsIdentifierPartES6(int(r)) {
			return false
		}
	}
	return true
}

// IsIdentifierES5 reports whether id is a valid ES5 identifier (name, not
// reserved).
func IsIdentifierES5(id string, strict bool) bool {
	return IsIdentifierNameES5(id) && !IsReservedWordES5(id, strict)
}

// IsIdentifierES6 reports whether id is a valid ES6 identifier.
func IsIdentifierES6(id string, strict bool) bool {
	return IsIdentifierNameES6(id) && !IsReservedWordES6(id, strict)
}
