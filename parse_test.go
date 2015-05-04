package haddoque

import (
	"fmt"
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

		fmt.Println(printIndentRoot(tr.root))
		equals(t, printIndentRoot(test.tree.root), printIndentRoot(tr.root))
	}
}
