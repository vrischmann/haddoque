package haddoque

import "testing"

// heavily based on text/template from the Go distribution

type lexTest struct {
	name  string
	input string
	items []lexeme
}

var (
	tEOF      = lexeme{tokEOF, 0, ""}
	tComma    = lexeme{tokComma, 0, ","}
	tWhere    = lexeme{tokWhere, 0, "where"}
	tEq       = lexeme{tokEq, 0, "=="}
	tNeq      = lexeme{tokNeq, 0, "!="}
	tNot      = lexeme{tokNot, 0, "!"}
	tGt       = lexeme{tokGt, 0, ">"}
	tGte      = lexeme{tokGte, 0, ">="}
	tLt       = lexeme{tokLt, 0, "<"}
	tLte      = lexeme{tokLte, 0, "<="}
	tAnd      = lexeme{tokAnd, 0, "and"}
	tOr       = lexeme{tokOr, 0, "or"}
	tIn       = lexeme{tokIn, 0, "in"}
	tContains = lexeme{tokContains, 0, "contains"}
	tLparen   = lexeme{tokLparen, 0, "("}
	tRparen   = lexeme{tokRparen, 0, ")"}
	tLbracket = lexeme{tokLbracket, 0, "["}
	tRbracket = lexeme{tokRbracket, 0, "]"}
)

var lexTests = []lexTest{
	{"empty", "", []lexeme{tEOF}},
	{"root field", ".", []lexeme{
		{tokField, 0, "."},
		tEOF,
	}},
	{"single field", ".name", []lexeme{
		{tokField, 0, ".name"},
		tEOF,
	}},
	{"nested field", ".data.id", []lexeme{
		{tokField, 0, ".data"},
		{tokField, 0, ".id"},
		tEOF,
	}},
	{"multiple fields", ".data.id, .name", []lexeme{
		{tokField, 0, ".data"},
		{tokField, 0, ".id"},
		tComma,
		{tokField, 0, ".name"},
		tEOF,
	}},
	{"with conditions", `.name where (.id == 1 and .age > 0.3) or .name != "foobar"`, []lexeme{
		{tokField, 0, ".name"},
		tWhere,
		tLparen,
		{tokField, 0, ".id"},
		tEq,
		{tokNumber, 0, "1"},
		tAnd,
		{tokField, 0, ".age"},
		tGt,
		{tokNumber, 0, "0.3"},
		tRparen,
		tOr,
		{tokField, 0, ".name"},
		tNeq,
		{tokString, 0, `"foobar"`},
		tEOF,
	}},
	{"with gte lte", ".id where .age >= 2 and .age <= 100", []lexeme{
		{tokField, 0, ".id"},
		tWhere,
		{tokField, 0, ".age"},
		tGte,
		{tokNumber, 0, "2"},
		tAnd,
		{tokField, 0, ".age"},
		tLte,
		{tokNumber, 0, "100"},
		tEOF,
	}},
	{"with in condition", ". where .age in [1, 2]", []lexeme{
		{tokField, 0, "."},
		tWhere,
		{tokField, 0, ".age"},
		tIn,
		tLbracket,
		{tokNumber, 0, "1"},
		tComma,
		{tokNumber, 0, "2"},
		tRbracket,
		tEOF,
	}},
	{"with contains condition", `. where .names contains "foobar"`, []lexeme{
		{tokField, 0, "."},
		tWhere,
		{tokField, 0, ".names"},
		tContains,
		{tokString, 0, `"foobar"`},
		tEOF,
	}},
	{"with contains list condition", `. where .names contains ["foobar", "barbaz"]`, []lexeme{
		{tokField, 0, "."},
		tWhere,
		{tokField, 0, ".names"},
		tContains,
		tLbracket,
		{tokString, 0, `"foobar"`},
		tComma,
		{tokString, 0, `"barbaz"`},
		tRbracket,
		tEOF,
	}},
	{"with not condition", `. where !( .age == 10 and .name == "foobar" )`, []lexeme{
		{tokField, 0, "."},
		tWhere,
		tNot,
		tLparen,
		{tokField, 0, ".age"},
		tEq,
		{tokNumber, 0, "10"},
		tAnd,
		{tokField, 0, ".name"},
		tEq,
		{tokString, 0, `"foobar"`},
		tRparen,
		tEOF,
	}},
}

func collect(t *lexTest) (items []lexeme) {
	l := newLexer(t.input)
	l.lex()
	for {
		item := l.nextLexeme()
		items = append(items, item)
		if item.tok == tokEOF || item.tok == tokError {
			break
		}
	}

	return
}

func equalLexemes(t testing.TB, i1, i2 []lexeme) {
	equals(t, len(i1), len(i2))
	for k := range i1 {
		equals(t, i1[k].tok, i2[k].tok)
		equals(t, i1[k].val, i2[k].val)
	}
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		equalLexemes(t, test.items, items)
	}
}
