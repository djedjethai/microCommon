package tests

import (
	// "encoding/json"
	"fmt"
	// "os"
	"reflect"
	"runtime"
	"testing"
)

// FailFunc is a function to call in case of failure
type FailFunc func(string, ...error)

// MaybeFail represents a fail function
var MaybeFail FailFunc

func InitFailFunc(t *testing.T) FailFunc {
	tester := t
	return func(msg string, errors ...error) {
		for _, err := range errors {
			if err != nil {
				pc := make([]uintptr, 1)
				runtime.Callers(2, pc)
				caller := runtime.FuncForPC(pc[0])
				_, line := caller.FileLine(caller.Entry())

				tester.Fatalf("%s:%d failed: %s %s", caller.Name(), line, msg, err)
			}
		}
	}
}

func Expect(actual, expected interface{}) error {
	if !reflect.DeepEqual(actual, expected) {
		return fmt.Errorf("expected: %v, Actual: %v", expected, actual)
	}

	return nil
}

// type testConf map[string]interface{}
// var testconf = make(testConf)
//
// // var srClient Client
// var maybeFail failFunc
//
// // getObject returns a child object of the root testConf
// func (tc testConf) getObject(name string) testConf {
// 	return tc[name].(map[string]interface{})
// }
//
// // getString returns a string representation of the value represented by key from the provided namespace
// // if the namespace is an empty string the root object will be searched.
// func (tc testConf) getString(key string) string {
// 	val, ok := tc[key]
// 	if ok {
// 		return val.(string)
// 	}
// 	return ""
// }
//
// // getInt returns an integer representation of the value represented by key from the provided namespace
// // If the namespace is an empty string the root object will be searched.
// func (tc testConf) getInt(key string) int {
// 	val, ok := tc[key]
// 	if ok {
// 		return val.(int)
// 	}
// 	return 0
// }
