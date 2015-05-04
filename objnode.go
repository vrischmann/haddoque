package haddoque

import (
	"errors"
	"strings"
)

var (
	errStopWalk = errors.New("haddoque: walk stopped")
)

type objNode struct {
	name   string
	value  interface{}
	fields []*objNode

	path        string
	cachedPaths map[string]struct{}
}

type walkObjNodeFn func(path string, node *objNode) error

func walkObjNode(path string, root *objNode, fn walkObjNodeFn) error {
	if err := fn(path, root); err != nil {
		return err
	}

	for _, v := range root.fields {
		p := makePath(path, v.name)
		if err := walkObjNode(p, v, fn); err != nil {
			return err
		}
	}

	return nil
}

func makePath(path, name string) string {
	return strings.Join([]string{path, name}, ".")
}

func newObjNode(obj interface{}) *objNode {
	v, ok := obj.(map[string]interface{})
	if !ok {
		return nil
	}

	return newObjNode1(&objNode{}, ".", v)
}

func newObjNode1(on *objNode, name string, obj interface{}) *objNode {
	on.name = name

	switch v := obj.(type) {
	case map[string]interface{}:
		for k, el := range v {
			newOn := newObjNode1(&objNode{}, k, el)
			on.fields = append(on.fields, newOn)
		}
	case interface{}:
		on.value = obj
	}

	return on
}

func (n *objNode) makeAllPaths() []string {
	var res []string
	walkObjNode("", n, func(path string, node *objNode) error {
		if path == "" {
			node.path = "."
		} else {
			node.path = path
		}
		res = append(res, path)
		return nil
	})
	return res
}

func (n *objNode) populateCachedPaths() {
	if len(n.cachedPaths) == 0 {
		n.cachedPaths = make(map[string]struct{})
		for _, p := range n.makeAllPaths() {
			n.cachedPaths[p] = struct{}{}
		}
		n.cachedPaths["."] = struct{}{}
	}
}

func (n *objNode) hasPath(path string) bool {
	n.populateCachedPaths()
	_, ok := n.cachedPaths[path]
	return ok
}

func (n *objNode) data() interface{} {
	if n.value != nil {
		return n.value
	}

	res := make(map[string]interface{})
	for _, f := range n.fields {
		res[f.name] = f.data()
	}

	return res
}

func (n *objNode) findSubNode(path string) *objNode {
	if n.path == path {
		return n
	}

	for _, f := range n.fields {
		if sn := f.findSubNode(path); sn != nil {
			return sn
		}
	}

	return nil
}

func (n *objNode) get(path string) interface{} {
	n.populateCachedPaths()
	if _, ok := n.cachedPaths[path]; !ok {
		return nil
	}

	sn := n.findSubNode(path)
	if sn == nil {
		return nil
	}

	return sn.data()
}

func (n *objNode) chainParts() []string {
	return strings.Split(n.path, ".")
}

func mergeNodes(nodes []*objNode) map[string]interface{} {
	if len(nodes) == 1 {
		return nodes[0].data().(map[string]interface{})
	}

	res := make(map[string]interface{})
	for _, n := range nodes {
		makeChainParts(res, n.chainParts(), n.data())
	}

	return res
}

func makeChainParts(m map[string]interface{}, parts []string, data interface{}) {
	current := m
	for i, p := range parts[1:] {
		if i+1 >= len(parts)-1 {
			current[p] = data
			break
		}
		current[p] = make(map[string]interface{})
		current = current[p].(map[string]interface{})
	}
}
