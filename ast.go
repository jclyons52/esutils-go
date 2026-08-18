package esutils

// AST-node predicates, mirroring esutils lib/ast.js. Each takes a minimal
// *Node; nil nodes answer false (matching JS null).

// IsExpression reports whether the node is an ESTree expression form.
func IsExpression(node *Node) bool {
	if node == nil {
		return false
	}
	switch node.Type {
	case "ArrayExpression", "AssignmentExpression", "BinaryExpression",
		"CallExpression", "ConditionalExpression", "FunctionExpression",
		"Identifier", "Literal", "LogicalExpression", "MemberExpression",
		"NewExpression", "ObjectExpression", "SequenceExpression",
		"ThisExpression", "UnaryExpression", "UpdateExpression":
		return true
	}
	return false
}

// IsIterationStatement reports whether the node is a looping statement.
func IsIterationStatement(node *Node) bool {
	if node == nil {
		return false
	}
	switch node.Type {
	case "DoWhileStatement", "ForInStatement", "ForStatement", "WhileStatement":
		return true
	}
	return false
}

// IsStatement reports whether the node is a statement form.
func IsStatement(node *Node) bool {
	if node == nil {
		return false
	}
	switch node.Type {
	case "BlockStatement", "BreakStatement", "ContinueStatement",
		"DebuggerStatement", "DoWhileStatement", "EmptyStatement",
		"ExpressionStatement", "ForInStatement", "ForStatement", "IfStatement",
		"LabeledStatement", "ReturnStatement", "SwitchStatement",
		"ThrowStatement", "TryStatement", "VariableDeclaration",
		"WhileStatement", "WithStatement":
		return true
	}
	return false
}

// IsSourceElement reports whether the node is a statement or a top-level
// function declaration.
func IsSourceElement(node *Node) bool {
	return node != nil && (IsStatement(node) || node.Type == "FunctionDeclaration")
}

// trailingStatement returns the statement that follows the given node in a
// control-flow tail position, or nil.
func trailingStatement(node *Node) *Node {
	switch node.Type {
	case "IfStatement":
		if node.Alternate != nil {
			return node.Alternate
		}
		return node.Consequent
	case "LabeledStatement", "ForStatement", "ForInStatement",
		"WhileStatement", "WithStatement":
		return node.Body
	}
	return nil
}

// IsProblematicIfStatement reports whether an IfStatement has an alternate
// whose chain nests an if-without-alternate in tail position — a shape that
// is ambiguous to parse (dangling else). Mirrors esutils ast.js.
func IsProblematicIfStatement(node *Node) bool {
	if node == nil || node.Type != "IfStatement" || node.Alternate == nil {
		return false
	}
	current := node.Consequent
	for current != nil {
		if current.Type == "IfStatement" && current.Alternate == nil {
			return true
		}
		current = trailingStatement(current)
	}
	return false
}
