package haddoque

import (
	"errors"
	"strings"
)

var (
	ErrNonExistingFields = errors.New("some requested fields do not exist")
	ErrInvalidObject     = errors.New("unable to use the provided object")
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(query string, obj interface{}) (interface{}, error) {
	lexer := newLexer(query)
	lexer.lex()
	tr := newTree(lexer)

	err := tr.parse()
	if err != nil {
		return nil, err
	}

	on := newObjNode(obj)
	if on == nil {
		return nil, ErrInvalidObject
	}

	// first check the fields exists
	if !checkFields(tr.root, on) {
		return nil, ErrNonExistingFields
	}

	if !evaluateWhere(tr.root, on) {
		return nil, nil
	}

	return getFields(tr.root, on)
}

func checkFields(root *seqNode, on *objNode) bool {
	for _, v := range root.nodes {
		if v.typ() != nodeChain {
			break
		}

		if !on.hasPath(v.(*chainNode).chain) {
			return false
		}
	}

	return true
}

func evaluateWhere(root *seqNode, on *objNode) bool {
	var wn *whereNode
	for _, v := range root.nodes {
		if v.typ() == nodeWhere {
			wn = v.(*whereNode)
			break
		}
	}

	// if no WHERE node it's still okay
	if wn == nil {
		return true
	}

	return evaluateCondition(wn.condition, on)
}

func evaluateCondition(cond node, on *objNode) bool {
	switch v := cond.(type) {
	case *andNode:
		return evaluateCondition(v.left, on) && evaluateCondition(v.right, on)
	case *orNode:
		return evaluateCondition(v.left, on) || evaluateCondition(v.right, on)
	case *inNode:
		if v.right.typ() != nodeSeq {
			return false
		}

		l := v.right.(*seqNode)
		for _, el := range l.nodes {
			if evaluateCondition(el, on) {
				return true
			}
		}

		return false
	case *containsNode:
	case *operationNode:
		return evaluateOperationNode(v, on)
	}

	return false
}

func evaluateOperationNode(n *operationNode, on *objNode) bool {
	// TODO(vincent): do we want to support something else as LHS ?
	if n.left.typ() != nodeChain {
		return false
	}

	lval := on.get(n.left.(*chainNode).chain) // we know for sure the object has that field at this point in time

	// TODO(vincent): does the parser allow non value nodes here ? need to check
	rval := getValue(n.right, on)
	if rval == nil {
		return false
	}

	switch n.operator {
	case tokLt: // <
		return evaluateLt(lval, rval)
	case tokLte: // <=
		return evaluateLte(lval, rval)
	case tokGt: // >
		return evaluateGt(lval, rval)
	case tokGte: // >=
		return evaluateGt(lval, rval)
	case tokEq: // ==
		return evaluateEq(lval, rval)
	case tokNeq: // !=
		return evaluateNeq(lval, rval)
	case tokNot: // !
		return false // TODO(vincent): unsupported, can we handle this ?
	}

	return false
}

// getFields selects the wanted fields from the objNode
func getFields(root *seqNode, on *objNode) (interface{}, error) {
	var nodes []*objNode

	// Beware: this is ugly code

	for _, v := range root.nodes {
		if v.typ() != nodeChain {
			break
		}

		chain := v.(*chainNode).chain
		nodes = append(nodes, on.findSubNode(chain))
	}

	return mergeNodes(nodes), nil
}

func getValue(n node, on *objNode) interface{} {
	switch v := n.(type) {
	case *chainNode:
		return on.get(v.chain)
	case *boolNode:
		return v.val
	case *textNode:
		return v.text
	case *numberNode:
		if v.isInt {
			return v.intVal
		}

		return v.floatVal
	default:
		return nil
	}
}

func evaluateLt(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv < rv
		}

		frv, ok := r.(float64)
		return ok && lv < int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv < rv
		}

		irv, ok := r.(int64)
		return ok && lv < float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv < rv
	default:
		return false
	}
}

func evaluateLte(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv <= rv
		}

		frv, ok := r.(float64)
		return ok && lv <= int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv <= rv
		}

		irv, ok := r.(int64)
		return ok && lv <= float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv <= rv
	default:
		return false
	}
}

func evaluateGt(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv > rv
		}

		frv, ok := r.(float64)
		return ok && lv > int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv > rv
		}

		irv, ok := r.(int64)
		return ok && lv > float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv > rv
	default:
		return false
	}
}

func evaluateGte(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv >= rv
		}

		frv, ok := r.(float64)
		return ok && lv >= int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv >= rv
		}

		irv, ok := r.(int64)
		return ok && lv >= float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv >= rv
	default:
		return false
	}
}

func evaluateEq(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv == rv
		}

		frv, ok := r.(float64)
		return ok && lv == int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv == rv
		}

		irv, ok := r.(int64)
		return ok && lv == float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv == rv
	default:
		return false
	}
}

func evaluateNeq(l, r interface{}) bool {
	switch lv := l.(type) {
	case int64:
		rv, ok := r.(int64)
		if ok {
			return lv != rv
		}

		frv, ok := r.(float64)
		return ok && lv != int64(frv)
	case float64:
		rv, ok := r.(float64)
		if ok {
			return lv != rv
		}

		irv, ok := r.(int64)
		return ok && lv != float64(irv)
	case string:
		rv, ok := r.(string)
		if !ok {
			return false
		}
		rv = strings.Trim(rv, `"`)
		return ok && lv != rv
	default:
		return false
	}
}
