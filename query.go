package haddoque

import "fmt"

/*

Syntax which is easier to parse

.name .data.id |
or(
    eq(.id, 1)
    and(
        neq(.name, "foobar")
        lt(.id, 100)
    )
)

Or I keep the current syntax but I'm having trouble finding a way to parse it

*/

type query interface {
	exec(node *objNode) (interface{}, error)
}

type queryFunc func(node *objNode) (interface{}, error)

func (f queryFunc) exec(node *objNode) (interface{}, error) {
	return f(node)
}

type expression interface {
	expr()
}

func (e *binaryExpression) expr() {}
func (e *valueExpression) expr()  {}
func (e *fieldExpression) expr()  {}

//go:generate stringer -type=operator
type operator int

const (
	opUnknown operator = iota
	opAnd
	opOr
	opLte
	opGte
	opLt
	opGt
	opEq
	opNeq
)

type fieldExpression struct {
	field string
}

type valueExpression struct {
	value interface{}
}

type binaryExpression struct {
	op    operator
	left  expression
	right expression
}

type filteringQuery struct {
	fields    []string
	condition expression
}

func (q *filteringQuery) exec(node *objNode) (interface{}, error) {
	for _, f := range q.fields {
		if !node.hasPath(f) {
			// TODO(vincent): better error
			return nil, fmt.Errorf("field %v does not exist in node", f)
		}
	}

	return evaluate(node, q.condition)
}

func evaluate(node *objNode, condition expression) (interface{}, error) {
	switch c := condition.(type) {
	case *valueExpression:
		return c.value, nil
	case *fieldExpression:
		return node.get(c.field), nil
	case *binaryExpression:
		return evaluateBinaryExpression(node, c)
	}

	return nil, nil
}

func evaluateBinaryExpression(node *objNode, c *binaryExpression) (res interface{}, err error) {
	var leftValue, rightValue interface{}

	if leftValue, err = evaluate(node, c.left); err != nil {
		return nil, err
	}
	if rightValue, err = evaluate(node, c.right); err != nil {
		return nil, err
	}

	switch c.op {
	case opAnd:
		res = isValueValid(leftValue) && isValueValid(rightValue)
	case opOr:
		res = isValueValid(leftValue) || isValueValid(rightValue)
	case opLte:
		res, err = evaluateLte(leftValue, rightValue)
	case opGte:
		res, err = evaluateGte(leftValue, rightValue)
	case opLt:
		res, err = evaluateLt(leftValue, rightValue)
	case opGt:
		res, err = evaluateGt(leftValue, rightValue)
	case opEq:
		res = leftValue == rightValue
	case opNeq:
		res = leftValue != rightValue
	default:
		err = fmt.Errorf("unknown operator %v", c.op)
	}

	return
}

func isValueValid(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	default:
		return v != nil
	}
}

func evaluateLte(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		r, ok := right.(int)
		if !ok {
			return nil, fmt.Errorf("right value of <= expression (%v) is not comparable to left value", right)
		}
		return l <= r, nil
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("right value of <= expression (%v) is not comparable to left value", right)
		}
		return l <= r, nil
	default:
		return nil, fmt.Errorf("unsupported left value type %T", l)
	}
}

func evaluateGte(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		r, ok := right.(int)
		if !ok {
			return nil, fmt.Errorf("right value of >= expression (%v) is not comparable to left value", right)
		}
		return l >= r, nil
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("right value of >= expression (%v) is not comparable to left value", right)
		}
		return l >= r, nil
	default:
		return nil, fmt.Errorf("unsupported left value type %T", l)
	}
}

func evaluateLt(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		r, ok := right.(int)
		if !ok {
			return nil, fmt.Errorf("right value of < expression (%v) is not comparable to left value", right)
		}
		return l < r, nil
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("right value of < expression (%v) is not comparable to left value", right)
		}
		return l < r, nil
	default:
		return nil, fmt.Errorf("unsupported left value type %T", l)
	}
}

func evaluateGt(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		r, ok := right.(int)
		if !ok {
			return nil, fmt.Errorf("right value of > expression (%v) is not comparable to left value", right)
		}
		return l > r, nil
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("right value of > expression (%v) is not comparable to left value", right)
		}
		return l > r, nil
	default:
		return nil, fmt.Errorf("unsupported left value type %T", l)
	}
}
