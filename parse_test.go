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
		root: &listNode{nodeType: nodeList},
	}},
	{"root field", ".", &tree{
		root: &listNode{nodeType: nodeList, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: "."},
		}},
	}},
	{"single field", ".name", &tree{
		root: &listNode{nodeType: nodeList, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".name"},
		}},
	}},
	{"chain", ".data.id", &tree{
		root: &listNode{nodeType: nodeList, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".data.id"},
		}},
	}},
	{"multiple chains", ".data.id, .name", &tree{
		root: &listNode{nodeType: nodeList, nodes: []node{
			&chainNode{nodeType: nodeChain, chain: ".data.id"},
			&chainNode{nodeType: nodeChain, chain: ".name"},
		}},
	}},
	{"with conditions", `.name where (.id == 1 and .age > 0.3) or .name != "foobar"`, &tree{
		root: &listNode{nodeType: nodeList, nodes: []node{
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
						left:     &chainNode{nodeType: nodeChain, chain: ".name"},
						right:    &textNode{nodeType: nodeText, text: "foobar"},
						operator: tokNeq,
					},
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
	for _, test := range parseTests {
		tr := parse(t, &test)
		fmt.Println(test.tree.root.String())
		fmt.Println(tr.root.String())
		// equals(t, test.tree.root.String(), tr.root.String())
	}
}
