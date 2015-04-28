package haddoque

import (
	"sort"
	"testing"
)

var (
	r = &objNode{
		fields: []*objNode{
			&objNode{name: "data", fields: []*objNode{
				&objNode{name: "id", value: 1},
				&objNode{name: "name", value: "Vincent"},
				&objNode{name: "platform", fields: []*objNode{
					&objNode{name: "type", value: "mobile"},
					&objNode{name: "value", value: "android"},
				}},
			}},
			&objNode{name: "locale", fields: []*objNode{
				&objNode{name: "language", value: "fr"},
				&objNode{name: "region", value: "FR"},
			}},
		},
	}
)

func TestObjNode(t *testing.T) {
	paths := r.makeAllPaths()
	sort.Strings(paths)

	exp := []string{
		"", ".data", ".data.id", ".data.name", ".data.platform",
		".data.platform.type", ".data.platform.value",
		".locale", ".locale.language", ".locale.region",
	}
	equals(t, exp, paths)
}

func TestObjNodeHasPath(t *testing.T) {
	equals(t, true, r.hasPath(".data"))
	equals(t, true, r.hasPath(".data.platform.type"))
	equals(t, false, r.hasPath(".foobar"))
}

func TestObjNodeGet(t *testing.T) {
	equals(t, 1, r.get(".data.id"))
	equals(t, nil, r.get(".data"))
	equals(t, "FR", r.get(".locale.region"))
}
