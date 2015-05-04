package haddoque

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// this is healivy based on the code from text/template

// tree is the representation of a parsed query
type tree struct {
	root  *seqNode
	lexer *lexer
	// buffer for peeking
	peekBuffer [2]lexeme
	peekCount  int
}

func newTree(lexer *lexer) *tree {
	return &tree{
		lexer: lexer,
	}
}

// nextLexeme fetches the next lexeme.
func (t *tree) nextLexeme() lexeme {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.peekBuffer[0] = t.lexer.nextLexeme()
	}
	return t.peekBuffer[t.peekCount]
}

// peek returns the next lexeme without consuming it.
func (t *tree) peek() lexeme {
	if t.peekCount > 0 {
		return t.peekBuffer[t.peekCount-1]
	}

	t.peekCount = 1
	t.peekBuffer[0] = t.lexer.nextLexeme()

	return t.peekBuffer[0]
}

// backup goes back one lexeme.
func (t *tree) backup() {
	t.peekCount++
}

// errorf panics with a formatted error.
func (t *tree) errorf(format string, args ...interface{}) {
	t.root = nil
	panic(fmt.Errorf(format, args...))
}

// recover catches panics and set the attached error to errp if it's not a runtime.Error
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

// parse starts parsing the query
func (t *tree) parse() (err error) {
	defer t.recover(&err)
	t.root = newSeqNode()

	p := t.peek()
	for ; p.tok != tokEOF; p = t.peek() {
		switch p.tok {
		case tokComma:
			t.nextLexeme()
		case tokField:
			n := t.parseChain(new([]string))
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

// parseChain parses a chain of fields
//
// NOTE(vincent): with how it's written right now, we need to pass a non empty slice at the first call.
// It'll need to be changed in the future.
func (t *tree) parseChain(fields *[]string) node {
	n := &chainNode{nodeType: nodeChain}
	defer func() {
		n.chain = strings.Join(*fields, "")
		fields = nil
	}()

	switch l := t.nextLexeme(); {
	case l.tok == tokField:
		*fields = append(*fields, l.val)
		return t.parseChain(fields)
	}
	t.backup()

	return n
}

// parseWhere parses a WHERE construct
func (t *tree) parseWhere() node {
	t.nextLexeme()

	n := &whereNode{nodeType: nodeWhere}
	n.condition = t.parseCondition()

	ty := n.condition.typ()
	if ty != nodeOr && ty != nodeAnd &&
		ty != nodeIn && ty != nodeContains &&
		ty != nodeOperation {
		t.errorf("unexpected condition")
	}

	return n
}

// parseCondition parses a condition
func (t *tree) parseCondition() node {
	var n node

	// Right now conditions are the last piece of a query, so parseCondition reads everything
	// until EOF.
	//
	// Because it's easier to handle, we require that every condition is enclosed in a pair of ().
	//
	// It's probably not a good implementation, but bear with me, I only just now started working with complex parsers.
	//
	// Here an overview of how it works:
	//  - a condition is just a node that can later be interpreted as "truthy".
	//    Obviously any expression which evaluates to a bool are valid,
	//    but single values like a number of string are valid too.
	//  - since we don't know in advance what kind of expression we have on our hand
	//    (remember, we use the grammer <lhs> <op> <rhs>), first time round
	//    we consider the expression a unary expression.
	//
	// Next time round, if we encounter an operator, we simply transform the unary expression
	// into a binary expression, and continue reading for the right hand side.
	for l := t.peek(); l.tok != tokEOF; l = t.peek() {
		switch {
		case l.tok == tokField:
			// TODO(vincent): error handling !
			if n == nil {
				n = t.parseChain(new([]string))
			} else if isBinaryExprNode(n) {
				t.setRightNode(n, t.parseChain(new([]string)))
			}
		case l.tok > tokLiteralsBegin && l.tok < tokLiteralsEnd:
			if n == nil {
				n = t.parseLiteral()
			} else if isBinaryExprNode(n) {
				t.setRightNode(n, t.parseLiteral())
			}
		case l.tok > tokOperatorsBegin && l.tok < tokOperatorsEnd:
			t.nextLexeme()
			n = &operationNode{
				nodeType: nodeOperation,
				left:     n,
				operator: l.tok,
			}
		case l.tok > tokKeywordsBegin && l.tok < tokKeywordsEnd:
			n = t.parseKeyword(n)
		case l.tok == tokLbracket:
			// TODO(vincent): error handling !
			if isBinaryExprNode(n) {
				t.setRightNode(n, t.parseLiteralSeq())
			}
		case l.tok == tokLparen:
			t.nextLexeme() // consume
			if n == nil {
				n = t.parseCondition()
			} else if isBinaryExprNode(n) {
				t.setRightNode(n, t.parseCondition())
			}
		case l.tok == tokRparen:
			t.nextLexeme() // consume
			return n
		default:
			t.errorf("unexpected token %v", l.tok)
		}
	}

	return n
}

func isBinaryExprNode(n node) bool {
	switch n.(type) {
	case *andNode, *orNode, *inNode, *containsNode, *operationNode:
		return true
	default:
		return false
	}
}

func (t *tree) setRightNode(n node, r node) {
	switch v := n.(type) {
	case *andNode:
		v.right = r
	case *orNode:
		v.right = r
	case *inNode:
		v.right = r
	case *containsNode:
		v.right = r
	case *operationNode:
		v.right = r
	default:
		t.errorf("node not a binary expr node")
	}
}

// parseLiteral parses a literal value
func (t *tree) parseLiteral() node {
	var n node
	switch l := t.nextLexeme(); {
	case l.tok == tokBool:
		val := l.val == "true"
		n = &boolNode{
			nodeType: nodeBool,
			val:      val,
		}
	case l.tok == tokChar, l.tok == tokString:
		n = &textNode{
			nodeType: nodeText,
			text:     l.val,
		}
	case l.tok == tokNumber:
		t.backup()
		n = t.parseNumber()
	}

	return n
}

// parseNumber parses a number value
func (t *tree) parseNumber() node {
	var err error
	l := t.nextLexeme()
	n := &numberNode{nodeType: nodeNumber}

	if strings.ContainsAny(l.val, "e.") {
		n.isFloat = true
		n.floatVal, err = strconv.ParseFloat(l.val, 64)
		if err != nil {
			t.errorf("bad number syntax. err=%v", err)
		}
	} else {
		n.isInt = true
		n.intVal, err = strconv.ParseInt(l.val, 10, 64)
		if err != nil {
			t.errorf("bad number syntax. err=%v", err)
		}
	}

	return n
}

// parseKeyword parses a keyword
func (t *tree) parseKeyword(left node) node {
	switch l := t.nextLexeme(); {
	case l.tok == tokAnd:
		return &andNode{
			nodeType: nodeAnd,
			left:     left,
		}
	case l.tok == tokOr:
		return &orNode{
			nodeType: nodeOr,
			left:     left,
		}
	case l.tok == tokIn:
		return &inNode{
			nodeType: nodeIn,
			left:     left,
		}
	case l.tok == tokContains:
		return &containsNode{
			nodeType: nodeContains,
			left:     left,
		}
	}

	return nil
}

// parseLiteralSeq parses a sequence of literal values only
func (t *tree) parseLiteralSeq() node {
	if t.peek().tok != tokLbracket {
		return nil
	}
	t.nextLexeme()

	n := &seqNode{nodeType: nodeSeq}

	for l := t.peek(); l.tok != tokEOF; l = t.peek() {
		switch {
		case l.tok == tokComma:
			t.nextLexeme() // consume
		case l.tok > tokLiteralsBegin && l.tok < tokLiteralsEnd:
			n.nodes = append(n.nodes, t.parseLiteral())
		case l.tok == tokRbracket:
			t.nextLexeme()
			return n
		}
	}

	return nil
}
