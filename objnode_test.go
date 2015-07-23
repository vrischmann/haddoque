package haddoque

import (
	"sort"
	"testing"
)

var (
	r = &objNode{
		fields: []*objNode{
			{name: "data", fields: []*objNode{
				{name: "id", value: 1},
				{name: "name", value: "Vincent"},
				{name: "platform", fields: []*objNode{
					{name: "type", value: "mobile"},
					{name: "value", value: "android"},
				}},
			}},
			{name: "locale", fields: []*objNode{
				{name: "language", value: "fr"},
				{name: "region", value: "FR"},
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
	equals(t, "FR", r.get(".locale.region"))
}

func TestNewObjNode(t *testing.T) {
	m := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   1,
			"name": "Vincent",
			"platform": map[string]interface{}{
				"type":  "mobile",
				"value": "android",
			},
		},
		"locale": map[string]interface{}{
			"language": "fr",
			"region":   "FR",
		},
		"shards": []int{1, 2, 3},
	}

	on := newObjNode(m)
	paths := on.makeAllPaths()

	sort.Strings(paths)

	exp := []string{
		"", ".data", ".data.id", ".data.name", ".data.platform",
		".data.platform.type", ".data.platform.value",
		".locale", ".locale.language", ".locale.region",
		".shards",
	}
	equals(t, exp, paths)

	equals(t, 1, on.get(".data.id"))
	equals(t, "Vincent", on.get(".data.name"))
	equals(t, "mobile", on.get(".data.platform.type"))
	equals(t, "android", on.get(".data.platform.value"))
	equals(t, "fr", on.get(".locale.language"))
	equals(t, "FR", on.get(".locale.region"))
	equals(t, []int{1, 2, 3}, on.get(".shards"))
}
