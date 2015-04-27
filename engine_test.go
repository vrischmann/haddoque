package haddoque

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
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
