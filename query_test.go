package haddoque

import "testing"

func TestExpression(t *testing.T) {
	// eq
	condition := &binaryExpression{
		op:    opEq,
		left:  &fieldExpression{".data.platform.type"},
		right: &valueExpression{"mobile"},
	}

	res, err := evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	// lte
	condition = &binaryExpression{
		op:    opLte,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{100},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	condition = &binaryExpression{
		op:    opLte,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{1},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	// gte
	condition = &binaryExpression{
		op:    opGte,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{-100},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	condition = &binaryExpression{
		op:    opGte,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{1},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	// lt
	condition = &binaryExpression{
		op:    opLt,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{1},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, false, res)

	condition = &binaryExpression{
		op:    opLt,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{2},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	// gt
	condition = &binaryExpression{
		op:    opLt,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{1},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, false, res)

	condition = &binaryExpression{
		op:    opLt,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{2},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)

	// neq
	condition = &binaryExpression{
		op:    opNeq,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{1},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, false, res)

	condition = &binaryExpression{
		op:    opNeq,
		left:  &fieldExpression{".data.id"},
		right: &valueExpression{2},
	}

	res, err = evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)
}

func TestSubExpressions(t *testing.T) {
	condition := &binaryExpression{
		op: opAnd,
		left: &binaryExpression{
			op:    opEq,
			left:  &fieldExpression{".data.platform.type"},
			right: &valueExpression{"mobile"},
		},
		right: &binaryExpression{
			op:    opEq,
			left:  &fieldExpression{".data.id"},
			right: &valueExpression{1},
		},
	}

	res, err := evaluate(r, condition)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
	equals(t, true, res)
}

func TestFilteringQuery(t *testing.T) {
	q := filteringQuery{
		fields: []string{".data.id", ".data.platform"},
		condition: &binaryExpression{
			op:    opEq,
			left:  &fieldExpression{".data.platform.type"},
			right: &valueExpression{"mobile"},
		},
	}

	res, err := q.exec(r)
	ok(t, err)
	assert(t, res != nil, "res should not be nil")
}
