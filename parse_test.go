package haddoque

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type parseTest struct {
	name  string
	input string
	tree  *tree
}

var parseTests = []parseTest{
	{"empty query", "", &tree{
		root: &seqNode{nodeType: nodeSeq},
	}},
	{"root field", ".", &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: "."},
		}},
	}},
	{"single field", ".name", &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".name"},
		}},
	}},
	{"chain", ".data.id", &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".data.id"},
		}},
	}},
	{"multiple chains", ".data.id, .name", &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".data.id"},
			&chainNode{nodeType: nodeChain, chain: ".name"},
		}},
	}},
	{"one condition", `.name where (.data.id == 1) and (.name != "vincent")`, &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".name"},
			&whereNode{
				nodeType: nodeWhere,
				condition: &andNode{
					nodeType: nodeAnd,
					left: &operationNode{
						nodeType: nodeOperation,
						left:     &chainNode{nodeType: nodeChain, chain: ".data.id"},
						right:    &numberNode{nodeType: nodeNumber, isInt: true, intVal: 1},
						operator: tokEq,
					},
					right: &operationNode{
						nodeType: nodeOperation,
						left:     &chainNode{nodeType: nodeChain, chain: ".name"},
						right:    &textNode{nodeType: nodeText, text: `"vincent"`},
						operator: tokEq,
					},
				},
			},
		}},
	}},
	{"with conditions", `.name where ( (.id == 1) and (.age > 0.3) ) or (.data.name != "foobar")`, &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".name"},
			&whereNode{
				nodeType: nodeWhere,
				condition: &orNode{
					nodeType: nodeOr,
					left: &andNode{
						nodeType: nodeAnd,
						left: &operationNode{
							nodeType: nodeOperation,
							left:     &chainNode{nodeType: nodeChain, chain: ".id"},
							right:    &numberNode{nodeType: nodeNumber, isInt: true, intVal: 1},
							operator: tokEq,
						},
						right: &operationNode{
							nodeType: nodeOperation,
							left:     &chainNode{nodeType: nodeChain, chain: ".age"},
							right:    &numberNode{nodeType: nodeNumber, isFloat: true, floatVal: 0.3},
							operator: tokGt,
						},
					},
					right: &operationNode{
						nodeType: nodeOperation,
						left:     &chainNode{nodeType: nodeChain, chain: ".data.name"},
						right:    &textNode{nodeType: nodeText, text: `"foobar"`},
						operator: tokNeq,
					},
				},
			},
		}},
	}},
	{"with in", `.name where ( .id in [1, 2, 3] )`, &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".name"},
			&whereNode{
				nodeType: nodeWhere,
				condition: &inNode{
					nodeType: nodeIn,
					left:     &chainNode{nodeType: nodeChain, chain: ".id"},
					right: &seqNode{nodeType: nodeSeq, nodes: []node{
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 1},
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 2},
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 3},
					}},
				},
			},
		}},
	}},
	{"with contains", `.id where ( .shards contains [1, 2, 3] )`, &tree{
		root: &seqNode{nodeType: nodeSeq, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".id"},
			&whereNode{
				nodeType: nodeWhere,
				condition: &containsNode{
					nodeType: nodeContains,
					left:     &chainNode{nodeType: nodeChain, chain: ".shards"},
					right: &seqNode{nodeType: nodeSeq, nodes: []node{
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 1},
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 2},
						&numberNode{nodeType: nodeNumber, isInt: true, intVal: 3},
					}},
				},
			},
		}},
	}},
}

func parse(t testing.TB, test *parseTest) *tree {
	l := newLexer(test.input)
	l.lex()
	tr := newTree(l)

	err := tr.parse()
	ok(t, err)

	return tr
}

func TestParse(t *testing.T) {
	for _, test := range parseTests[len(parseTests)-1:] {
		tr := parse(t, &test)
		equals(t, printIndentRoot(test.tree.root), printIndentRoot(tr.root))
	}
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
