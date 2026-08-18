// Package esutils is a Go port of the estools/esutils library: utility
// predicates for ECMAScript AST nodes, code characters, and keywords.
package esutils

// Node is the subset of an ESTree AST node that esutils inspects. It accepts
// the AST shapes produced by ESTree-parsers; only the fields below are read.
type Node struct {
	Type       string `json:"type"`
	Consequent *Node  `json:"consequent,omitempty"`
	Alternate  *Node  `json:"alternate,omitempty"`
	Body       *Node  `json:"body,omitempty"`
}
