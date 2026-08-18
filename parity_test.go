package esutils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Entry is one parity case: an op plus its payload. Both this test and the
// JS driver (node + real esutils) compute a result for each Entry so the two
// can be compared 1:1.
type Entry struct {
	Op     string `json:"op"`
	Node   *Node  `json:"node,omitempty"`
	S      string `json:"s"`
	Ch     int    `json:"ch,omitempty"`
	Strict bool   `json:"strict,omitempty"`
}

// buildCorpus returns a representative, exhaustive-of-logic set of cases.
func buildCorpus() []Entry {
	var es []Entry
	nodeOps := []string{"isExpression", "isStatement", "isIterationStatement",
		"isSourceElement", "isProblematicIfStatement", "trailingStatement"}
	// The four null-safe predicates (JS null-checks them); isProblematicIf and
	// trailingStatement throw on null in the original, so a nil node is only
	// tested against the null-safe ones.
	nullSafeOps := nodeOps[:4]
	addNode := func(n *Node, ops []string) {
		for _, op := range ops {
			es = append(es, Entry{Op: op, Node: n})
		}
	}
	ifN := func(typ string, alt, cons *Node) *Node { return &Node{Type: typ, Alternate: alt, Consequent: cons} }
	bodyN := func(typ string, body *Node) *Node { return &Node{Type: typ, Body: body} }
	nodes := []*Node{
		{Type: "IfStatement", Consequent: &Node{Type: "ExpressionStatement"}},
		ifN("IfStatement", &Node{Type: "EmptyStatement"}, &Node{Type: "ExpressionStatement"}),
		ifN("IfStatement", nil, &Node{Type: "IfStatement", Alternate: &Node{Type: "EmptyStatement"}, Consequent: &Node{Type: "EmptyStatement"}}),
		{Type: "WhileStatement", Body: &Node{Type: "BlockStatement"}},
		bodyN("ForStatement", &Node{Type: "ExpressionStatement"}),
		{Type: "Identifier"},
		{Type: "Literal", Consequent: &Node{Type: "FunctionExpression"}},
		{Type: "FunctionDeclaration"}, {Type: "VariableDeclaration"}, {Type: "CallExpression"},
		{Type: "ReturnStatement"}, {Type: "BreakStatement"}, {Type: "ArrowFunctionExpression"},
		{Type: "WithStatement", Body: &Node{Type: "IfStatement"}},
	}
	for _, n := range nodes {
		addNode(n, nodeOps)
	}
	addNode(nil, nullSafeOps)

	words := []string{"if", "in", "do", "var", "for", "new", "try", "this", "else",
		"case", "void", "with", "enum", "while", "break", "catch", "throw", "const",
		"yield", "class", "super", "return", "typeof", "delete", "switch", "export",
		"import", "default", "finally", "extends", "function", "continue", "debugger",
		"instanceof", "implements", "interface", "package", "private", "protected",
		"public", "static", "let", "null", "true", "false", "eval", "arguments",
		"foo", "bar", "_", "$", "x1", "é", "π", "𐐀", "var1", "2x", "in", "fo o", "await", "async"}
	for _, w := range words {
		for _, strict := range []bool{false, true} {
			for _, op := range []string{"isKeywordES5", "isKeywordES6", "isReservedWordES5",
				"isReservedWordES6", "isRestrictedWord", "isIdentifierES5", "isIdentifierES6",
				"isIdentifierNameES5", "isIdentifierNameES6"} {
				es = append(es, Entry{Op: op, S: w, Strict: strict})
			}
		}
	}
	names := []string{"", "$", "_", "a", "a1", "1a", "hello", "class", "var",
		"παντελής", "日本語", "𐐀", "𠮷a", "a b", "a\x00b", "eval", "yield", "let"}
	for _, n := range names {
		es = append(es, Entry{Op: "isIdentifierNameES5", S: n})
		es = append(es, Entry{Op: "isIdentifierNameES6", S: n})
	}

	// character predicates: dense sweep of the common planes + supplementary samples
	toCharOps := []string{"isDecimalDigit", "isHexDigit", "isOctalDigit", "isWhiteSpace",
		"isLineTerminator", "isIdentifierStartES5", "isIdentifierPartES5",
		"isIdentifierStartES6", "isIdentifierPartES6"}
	for ch := 0; ch <= 0x2FFF; ch++ {
		for _, op := range toCharOps {
			es = append(es, Entry{Op: op, Ch: ch})
		}
	}
	for _, ch := range []int{0x10000, 0x10041, 0x1009D, 0x10140, 0x10145, 0x10400,
		0x1049D, 0x13400, 0x1D400, 0x1F1E6, 0x2A700, 0x10FFFF, 0x1308E, 0x2FFFF, 0xE0000} {
		for _, op := range toCharOps {
			es = append(es, Entry{Op: op, Ch: ch})
		}
	}
	return es
}

// goResult computes the Go port's answer for an entry.
func goResult(e Entry) any {
	switch e.Op {
	case "isExpression":
		return IsExpression(e.Node)
	case "isStatement":
		return IsStatement(e.Node)
	case "isIterationStatement":
		return IsIterationStatement(e.Node)
	case "isSourceElement":
		return IsSourceElement(e.Node)
	case "isProblematicIfStatement":
		return IsProblematicIfStatement(e.Node)
	case "trailingStatement":
		t := trailingStatement(e.Node)
		if t == nil {
			return nil
		}
		return t.Type
	case "isKeywordES5":
		return IsKeywordES5(e.S, e.Strict)
	case "isKeywordES6":
		return IsKeywordES6(e.S, e.Strict)
	case "isReservedWordES5":
		return IsReservedWordES5(e.S, e.Strict)
	case "isReservedWordES6":
		return IsReservedWordES6(e.S, e.Strict)
	case "isRestrictedWord":
		return IsRestrictedWord(e.S)
	case "isIdentifierES5":
		return IsIdentifierES5(e.S, e.Strict)
	case "isIdentifierES6":
		return IsIdentifierES6(e.S, e.Strict)
	case "isIdentifierNameES5":
		return IsIdentifierNameES5(e.S)
	case "isIdentifierNameES6":
		return IsIdentifierNameES6(e.S)
	case "isDecimalDigit":
		return IsDecimalDigit(e.Ch)
	case "isHexDigit":
		return IsHexDigit(e.Ch)
	case "isOctalDigit":
		return IsOctalDigit(e.Ch)
	case "isWhiteSpace":
		return IsWhiteSpace(e.Ch)
	case "isLineTerminator":
		return IsLineTerminator(e.Ch)
	case "isIdentifierStartES5":
		return IsIdentifierStartES5(e.Ch)
	case "isIdentifierPartES5":
		return IsIdentifierPartES5(e.Ch)
	case "isIdentifierStartES6":
		return IsIdentifierStartES6(e.Ch)
	case "isIdentifierPartES6":
		return IsIdentifierPartES6(e.Ch)
	}
	return "UNKNOWN_OP:" + e.Op
}

const jsDriver = `'use strict';
const fs = require('fs');
const log = require('util').log;
const corpus = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));
const code = require(process.env.ESUTILS_ORIG + '/lib/code.js');
const kw = require(process.env.ESUTILS_ORIG + '/lib/keyword.js');
const ast = require(process.env.ESUTILS_ORIG + '/lib/ast.js');
function r(e) {
  switch (e.op) {
    case 'isExpression': return ast.isExpression(e.node);
    case 'isStatement': return ast.isStatement(e.node);
    case 'isIterationStatement': return ast.isIterationStatement(e.node);
    case 'isSourceElement': return ast.isSourceElement(e.node);
    case 'isProblematicIfStatement': return ast.isProblematicIfStatement(e.node);
    case 'trailingStatement': { const t = ast.trailingStatement(e.node); return t ? t.type : null; }
    case 'isKeywordES5': return kw.isKeywordES5(e.s, e.strict);
    case 'isKeywordES6': return kw.isKeywordES6(e.s, e.strict);
    case 'isReservedWordES5': return kw.isReservedWordES5(e.s, e.strict);
    case 'isReservedWordES6': return kw.isReservedWordES6(e.s, e.strict);
    case 'isRestrictedWord': return kw.isRestrictedWord(e.s);
    case 'isIdentifierES5': return kw.isIdentifierES5(e.s, e.strict);
    case 'isIdentifierES6': return kw.isIdentifierES6(e.s, e.strict);
    case 'isIdentifierNameES5': return kw.isIdentifierNameES5(e.s);
    case 'isIdentifierNameES6': return kw.isIdentifierNameES6(e.s);
    case 'isDecimalDigit': return code.isDecimalDigit(e.ch);
    case 'isHexDigit': return code.isHexDigit(e.ch);
    case 'isOctalDigit': return code.isOctalDigit(e.ch);
    case 'isWhiteSpace': return code.isWhiteSpace(e.ch);
    case 'isLineTerminator': return code.isLineTerminator(e.ch);
    case 'isIdentifierStartES5': return code.isIdentifierStartES5(e.ch);
    case 'isIdentifierPartES5': return code.isIdentifierPartES5(e.ch);
    case 'isIdentifierStartES6': return code.isIdentifierStartES6(e.ch);
    case 'isIdentifierPartES6': return code.isIdentifierPartES6(e.ch);
  }
  return 'UNKNOWN_OP:' + e.op;
}
const out = [];
for (const e of corpus) out.push(r(e));
fs.writeFileSync(process.argv[3], JSON.stringify(out));
`

// TestParity runs the Go port against the real JS esutils and requires the
// outputs to agree for every corpus case. Skips when node is unavailable.
func TestParity(t *testing.T) {
	nodeBin, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not available; skipping JS parity")
	}
	orig := filepath.Join("original")
	if st, err := os.Stat(filepath.Join(orig, "lib", "code.js")); err != nil || st.IsDir() {
		t.Fatalf("original esutils not found at original/lib/code.js: %v", err)
	}
	corpus := buildCorpus()

	dir := t.TempDir()
	corpusJSON, _ := json.Marshal(corpus)
	if err := os.WriteFile(filepath.Join(dir, "corpus.json"), corpusJSON, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "driver.js"), []byte(jsDriver), 0o644); err != nil {
		t.Fatal(err)
	}
	origAbs, _ := filepath.Abs(orig)

	cmd := exec.Command(nodeBin, "driver.js", "corpus.json", "result.json")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ESUTILS_ORIG="+origAbs)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("node driver failed: %v\n%s", err, out)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	var jsResults []any
	if err := json.Unmarshal(raw, &jsResults); err != nil {
		t.Fatalf("parse node result: %v", err)
	}
	if len(jsResults) != len(corpus) {
		t.Fatalf("result length %d != corpus %d", len(jsResults), len(corpus))
	}

	mismatch := 0
	limit := 20
	for i, e := range corpus {
		got := goResult(e)
		want := stringify(jsResults[i])
		if stringify(got) != want {
			mismatch++
			if mismatch <= limit {
				t.Errorf("op=%-27s payload=%s go=%v js=%v", e.Op, key(e), got, jsResults[i])
			}
		}
	}
	t.Logf("parity: %d cases, %d mismatches", len(corpus), mismatch)
	if mismatch > 0 {
		t.Fatalf("parity mismatch (%d/%d) vs real esutils", mismatch, len(corpus))
	}
}

func stringify(v any) string {
	switch t := v.(type) {
	case bool:
		if t {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", t)
	}
}

func key(e Entry) string {
	var b strings.Builder
	b.WriteString(e.Op)
	if e.Node != nil {
		b.WriteString(":" + e.Node.Type)
	}
	if e.S != "" {
		b.WriteString(":" + strq(e.S))
	}
	if e.Ch != 0 {
		fmt.Fprintf(&b, ":%d", e.Ch)
	}
	if e.Strict {
		b.WriteString(":strict")
	}
	return b.String()
}

func strq(s string) string {
	r := s
	if len(r) > 12 {
		r = r[:12] + ".."
	}
	return r
}
