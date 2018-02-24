package util

import (
	"encoding/json"
	"testing"
)

type Foo struct {
	A int `json:"a"`
	B int `json:"b"`
}

// null struct
type Bar struct {
}

type Fee struct {
	A string `json:"aString"`
	B []int  `json:"arrayOfInts"`
}

var unmarshallTests []struct {
	in         string
	testStruct interface{}
	shouldErr  bool
	structOut  string
	unknownOut string
} = []struct {
	in         string
	testStruct interface{}
	shouldErr  bool
	structOut  string
	unknownOut string
}{
	// no json in
	{``, &Foo{}, true, "", ""},
	// incomplete json
	{`{"a":1, "b":2, "x"`, &Foo{}, true, "", ""},
	// both known and unknown
	{`{"a":1, "b":2, "x":3, "y":[6,7,8]}`, &Foo{}, false, `{"a":1,"b":2}`, `{"x":3,"y":[6,7,8]}`},
	// no unknown
	{`{"a":1, "b":2}`, &Foo{}, false, `{"a":1,"b":2}`, `{}`},
	// no known
	{`{"x":3, "y":[6,7,8]}`, &Foo{}, false, `{"a":0,"b":0}`, `{"x":3,"y":[6,7,8]}`},
	// empty struct
	{`{"a":1, "b":2, "x":3, "y":[6,7,8]}`, &Bar{}, false, `{}`, `{"a":1,"b":2,"x":3,"y":[6,7,8]}`},
	// with an array
	{`{"aString":"this is a string", "arrayOfInts":[6,7,8]}`, &Fee{}, false, `{"aString":"this is a string","arrayOfInts":[6,7,8]}`, `{}`},
}

func TestUnmarshalJsonObjAndMap(t *testing.T) {

	for index, tt := range unmarshallTests {
		unknown := map[string]json.RawMessage{}
		err := UnmarshalJson([]byte(tt.in), tt.testStruct, unknown)
		if (err == nil) && tt.shouldErr {
			t.Errorf("Test Index: %v: Unmarshal should of errored. returned nil", index)
		}
		if err == nil {
			actualUnkown, _ := json.Marshal(unknown)
			actualStruct, _ := json.Marshal(tt.testStruct)
			if (tt.unknownOut != string(actualUnkown)) || (tt.structOut != string(actualStruct)) {
				t.Errorf("Test Index: %v: Expected:\nIn Struct:%s\n Got:%s\nIn Unkown: %s\n Got %s\n", index, tt.structOut, actualStruct, tt.unknownOut, actualUnkown)
			}

		}

	}
}

var marshallTests []struct {
	in          string
	testStruct  interface{}
	shouldErr   bool
	expectedOut string
} = []struct {
	in          string
	testStruct  interface{}
	shouldErr   bool
	expectedOut string
}{
	// both known and unknown
	{`{"a":1, "b":2, "x":3, "y":[6,7,8]}`, &Foo{}, false, `{"x":3,"y":[6,7,8],"a":1,"b":2}`},
	// none
	{`{}`, &Foo{}, false, `{"a":0,"b":0}`},
	// no known
	{`{"x":3, "y":[6,7,8]}`, &Foo{}, false, `{"x":3,"y":[6,7,8],"a":0,"b":0}`},
	// No unknown
	{`{"a":1, "b":2}`, &Foo{}, false, `{"a":1,"b":2}`},
}

func TestMarshalJsonObjAndMap(t *testing.T) {

	for index, tt := range marshallTests {
		unknown := map[string]json.RawMessage{}
		err := UnmarshalJson([]byte(tt.in), tt.testStruct, unknown)
		if err != nil {
			t.Errorf("Test Index: %v: Umarshal Errored %v", index, err)
		}
		actualOutput, err := MarshalJson(tt.testStruct, unknown)
		if (err == nil) && tt.shouldErr {
			t.Errorf("Test Index: %v: marshal should of errored. returned nil", index)
		}
		if err == nil {
			if tt.expectedOut != string(actualOutput) {
				t.Errorf("Marshal Test Index: %v: Expected:%s\n Got:%s\n", index, tt.expectedOut, actualOutput)
			}
		}
	}
}
