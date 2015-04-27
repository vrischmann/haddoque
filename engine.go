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

func (n *objNode) makeAllPaths() []string {
	var res []string
	walkObjNode("", n, func(path string, node *objNode) error {
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
	}
}

func (n *objNode) hasPath(path string) bool {
	n.populateCachedPaths()
	_, ok := n.cachedPaths[path]
	return ok
}

func (n *objNode) get(path string) interface{} {
	n.populateCachedPaths()
	if _, ok := n.cachedPaths[path]; !ok {
		return nil
	}

	var res interface{}
	walkObjNode("", n, func(p string, node *objNode) error {
		if p == path {
			res = node.value
			return errStopWalk
		}

		return nil
	})
	return res
}

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(query string, obj interface{}) (interface{}, error) {
	return nil, nil
}
