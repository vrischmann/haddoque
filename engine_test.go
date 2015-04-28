package haddoque

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type engineTestData struct {
	input    interface{}
	expected interface{}
	query    string
}

type engineTest struct {
	file string
	data engineTestData
}

func readTest(t *testing.T, path string, input interface{}, query *string, expected interface{}) {
	data, err := ioutil.ReadFile("testdata/" + path)
	ok(t, err)

	tokens := bytes.Split(data, []byte("---"))
	tokens = tokens[1:]
	for _, v := range tokens {
		v = bytes.TrimSpace(v)
	}
	equals(t, 3, len(tokens))

	err = json.Unmarshal(tokens[0], input)
	ok(t, err)

	*query = string(tokens[1])

	err = json.Unmarshal(tokens[2], expected)
	ok(t, err)
}

var tests = []engineTest{
	engineTest{
		file: "1_simple_query.txt",
	},
}

func TestEngine(t *testing.T) {
	for _, test := range tests {
		readTest(t, test.file, &test.data.input, &test.data.query, &test.data.expected)

		engine := NewEngine()
		res, err := engine.Run(test.data.query, test.data.input)
		ok(t, err)
		equals(t, test.data.expected, res)
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
