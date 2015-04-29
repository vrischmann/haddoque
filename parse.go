package haddoque

import (
	"fmt"
	"runtime"
	"strings"
)

//go:generate stringer -type=nodeType
type nodeType int

const (
	nodeChain nodeType = iota
	nodeList
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

type node interface {
	typ() nodeType
	String() string
	// TODO(vincent): position ?
}

type listNode struct {
	nodeType
	nodes []node
}

func (l *listNode) String() string {
	return fmt.Sprintf("listNode{%v}", l.nodes)
}

func newListNode() *listNode {
	return &listNode{nodeType: nodeList}
}

type chainNode struct {
	nodeType
	chain string
}

func (c *chainNode) String() string {
	return fmt.Sprintf("chainNode{%s}", c.chain)
}

type boolNode struct {
	nodeType
	val bool
}

type textNode struct {
	nodeType
	text string
}

type numberNode struct {
	nodeType
	isInt    bool
	isFloat  bool
	intVal   int64
	floatVal float64
}

type whereNode struct {
	nodeType
	condition node
}

func (n *whereNode) String() string {
	return fmt.Sprintf("whereNode{%v}", n.condition)
}

type andNode struct {
	nodeType
	left  node
	right node
}

func (n *andNode) String() string {
	return fmt.Sprintf("andNode{left: %v, right: %v}", n.left, n.right)
}

type orNode struct {
	nodeType
	left  node
	right node
}

func (n *orNode) String() string {
	return fmt.Sprintf("orNode{left: %v, right: %v}", n.left, n.right)
}

type operationNode struct {
	nodeType
	left     node
	right    node
	operator token
}

func (n *operationNode) String() string {
	return fmt.Sprintf("operationNode{left: %v, right: %v, op: %s}", n.left, n.right, n.operator)
}

type tree struct {
	root       *listNode
	lexer      *lexer
	peekBuffer [2]lexeme
	peekCount  int
}

func newTree(lexer *lexer) *tree {
	return &tree{
		lexer: lexer,
	}
}

func walkTree1(root node, fn walkTreeFunc) error {
	if err := fn(root); err != nil {
		return err
	}

	switch n := root.(type) {
	case *listNode:
		for _, el := range n.nodes {
			if err := walkTree1(el, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkTree(tree *tree, fn walkTreeFunc) error {
	return walkTree1(tree.root, fn)
}

type walkTreeFunc func(n node) error

func (t *tree) nextLexeme() lexeme {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.peekBuffer[0] = t.lexer.nextLexeme()
	}
	return t.peekBuffer[t.peekCount]
}

func (t *tree) peek() lexeme {
	if t.peekCount > 0 {
		return t.peekBuffer[t.peekCount-1]
	}

	t.peekCount = 1
	t.peekBuffer[0] = t.lexer.nextLexeme()

	return t.peekBuffer[0]
}

func (t *tree) backup() {
	t.peekCount++
}

func (t *tree) errorf(format string, args ...interface{}) {
	t.root = nil
	panic(fmt.Errorf(format, args...))
}

func (t *tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}

		*errp = e.(error)
	}

	return
}

func (t *tree) parse() (err error) {
	defer t.recover(&err)
	t.root = newListNode()

	p := t.peek()
	for ; p.tok != tokEOF; p = t.peek() {
		switch p.tok {
		case tokComma:
			t.nextLexeme()
		case tokField:
			n := t.parseField(new([]string))
			t.root.nodes = append(t.root.nodes, n)
		case tokWhere:
			n := t.parseWhere()
			t.root.nodes = append(t.root.nodes, n)
		default:
			t.nextLexeme()
		}
	}

	return nil
}

func (t *tree) parseField(fields *[]string) node {
	n := &chainNode{nodeType: nodeChain}
	defer func() {
		n.chain = strings.Join(*fields, "")
		fields = nil
	}()

	switch l := t.nextLexeme(); {
	case l.tok == tokField:
		*fields = append(*fields, l.val)
		return t.parseField(fields)
	}
	t.backup()

	return n
}

func (t *tree) parseWhere() node {
	t.nextLexeme()

	n := &whereNode{nodeType: nodeWhere}
	n.condition = nil // TODO(vincent): how do I do this ?

	ty := n.condition.typ()
	if ty != nodeOr && ty != nodeAnd && ty != nodeOperation {
		t.errorf("unexpected condition")
	}

	return n
}
