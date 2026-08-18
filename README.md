# esutils-go

Go port of the npm package [`esutils`](https://github.com/estools/esutils) —
utility predicates for ECMAScript AST nodes, code characters, and keywords.

| | |
|---|---|
| Original size | ~477 LOC (lib/) |
| Public API | `IsExpression`, `IsStatement`, `IsIterationStatement`, `IsSourceElement`, `IsProblematicIfStatement`, `IsDecimalDigit`, `IsHexDigit`, `IsOctalDigit`, `IsWhiteSpace`, `IsLineTerminator`, `IsIdentifierStart/PartES5/ES6`, `IsKeywordES5/6`, `IsReservedWordES5/6`, `IsRestrictedWord`, `IsIdentifierNameES5/6`, `IsIdentifierES5/6` |
| Needed by | doctrine, eslint (etc.) in the uplift target |
| Module | `github.com/jclyons52/esutils-go` |

## Files
- `ast.go` — AST-node predicates over a minimal `Node{Type, Consequent, Alternate, Body}`.
- `code.go` — digit / whitespace / line-terminator checks + ASCII identifier tables.
- `code_tables.go` — **generated**: the four Unicode identifier start/part sets
  (ES5/ES6) as sorted code-point ranges with binary-search membership. Derived
  directly from the real JS module (not hand-transcribed from its giant regexes).
- `keyword.go` — ECMAScript keyword / reserved-word / identifier-name checks.
- `original/` — vendored original source (reference + parity oracle).
- `parity_test.go` — shells out to `node`, runs the real esutils over a
  111,949-case corpus, and requires identical results.

## Parity
```sh
go test ./...
# PARITY PASS: 111,949 cases, 0 mismatches (requires node)
```
