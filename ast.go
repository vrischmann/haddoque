package haddoque

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//go:generate stringer -type=nodeType
type nodeType int

const (
	nodeChain nodeType = iota
	nodeSeq
	nodeBool
	nodeText
	nodeNumber
	nodeWhere
	nodeAnd
	nodeOr
	nodeIn
	nodeContains
	nodeOperation
)

func (t nodeType) typ() nodeType {
	return t
}

// node is an element in the parse tree
type node interface {
	typ() nodeType
	String() string
	// TODO(vincent): position ?
}

// seqNode represents a sequence of values
type seqNode struct {
	nodeType
	nodes []node
}

func (l *seqNode) String() string {
	return "listNode"
}

func newSeqNode() *seqNode {
	return &seqNode{nodeType: nodeSeq}
}

// chainNode represents a chain of fields
type chainNode struct {
	nodeType
	chain string
}

func (c *chainNode) String() string {
	return fmt.Sprintf("chainNode{%s}", c.chain)
}

// boolNode represents a boolean value - true or false
type boolNode struct {
	nodeType
	val bool
}

func (n *boolNode) String() string {
	return fmt.Sprintf("boolNode{%v}", n.val)
}

// textNode represents a text value - string or char
type textNode struct {
	nodeType
	text string
}

func (n *textNode) String() string {
	return fmt.Sprintf("textNode{%s}", n.text)
}

// numberNode represents a number value - float or int
type numberNode struct {
	nodeType
	isInt    bool
	isFloat  bool
	intVal   int64
	floatVal float64
}

func (n *numberNode) String() string {
	if n.isInt {
		return fmt.Sprintf("nodeNumber{int: %d}", n.intVal)
	}
	return fmt.Sprintf("nodeNumber{float: %0.4f}", n.floatVal)
}

// whereNode represents a WHERE construct
type whereNode struct {
	nodeType
	condition node
}

func (n *whereNode) String() string {
	return "whereNode"
}

// andNode represents a binary AND expression
type andNode struct {
	nodeType
	left  node
	right node
}

func (n *andNode) String() string {
	return "andNode"
}

// orNode represents a binary OR expression
type orNode struct {
	nodeType
	left  node
	right node
}

func (n *orNode) String() string {
	return "orNode"
}

// operationNode represents a binary expression - >, <=, == etc
type operationNode struct {
	nodeType
	left     node
	right    node
	operator token
}

func (n *operationNode) String() string {
	return fmt.Sprintf("operationNode{%s}", n.operator)
}

// inNode represents a binary IN expression
type inNode struct {
	nodeType
	left  node
	right node
}

func (n *inNode) String() string {
	return fmt.Sprintf("inNode")
}

// containsNode represents a binary CONTAINS expression
type containsNode struct {
	nodeType
	left  node
	right node
}

func (n *containsNode) String() string {
	return fmt.Sprintf("containsNode")
}

// printIndentRoot prints an indented representation of the root
func printIndentRoot(root *seqNode) string {
	var buf bytes.Buffer
	printIndent(&buf, root, 0)

	return buf.String()
}

func printIndent(w io.Writer, root node, indent int) {
	if root == nil {
		return
	}
	fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", indent), root.String())

	switch v := root.(type) {
	case *seqNode:
		for _, el := range v.nodes {
			printIndent(w, el, indent+1)
		}
	case *whereNode:
		printIndent(w, v.condition, indent+1)
	case *andNode:
		printIndent(w, v.left, indent+1)
		printIndent(w, v.right, indent+1)
	case *orNode:
		printIndent(w, v.left, indent+1)
		printIndent(w, v.right, indent+1)
	case *inNode:
		printIndent(w, v.left, indent+1)
		printIndent(w, v.right, indent+1)
	case *containsNode:
		printIndent(w, v.left, indent+1)
		printIndent(w, v.right, indent+1)
	case *operationNode:
		printIndent(w, v.left, indent+1)
		printIndent(w, v.right, indent+1)
	}
}
