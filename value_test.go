package provider_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/klaidliadon/provider"
)

var (
	aVal = A{Str: "banana", Int: 42}
	bVal = B{Str: "banana", Int: 42}

	aJSON = []byte(`{"str":{"provider":"static","value":"banana"},"int":{"provider":"static","value":42}}`)
	bJSON = []byte(`{"provider":"static","value":{"str":"banana","int":42}}`)
)

type testCase struct {
	Value interface{}
	Bytes []byte
	Dst   interface{}
}

func TestMarshal(t *testing.T) {
	for _, tc := range []testCase{
		{&aVal, aJSON, nil},
		{&bVal, bJSON, nil},
	} {
		b, err := json.Marshal(tc.Value)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(b, tc.Bytes) {
			t.Fatalf("Expected:\n%s\nGot:\n%s", tc.Bytes, b)
		}
		t.Logf("%s", b)
	}
}

func TestUnmarshal(t *testing.T) {
	for _, tc := range []testCase{{&aVal, aJSON, &A{}}, {&bVal, bJSON, &B{}}} {
		if err := json.Unmarshal(tc.Bytes, tc.Dst); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tc.Dst, tc.Value) {
			t.Fatalf("Expected:\n%s\nGot:\n%s", tc.Value, tc.Dst)
		}
		t.Logf("%+v", tc.Dst)
	}
}

// A wraps each field
type A struct {
	Str string `json:"str,omitempty"`
	Int int    `json:"int,omitempty"`
}

func (a *A) toFields() interface{} {
	type A struct {
		Str *provider.Field `json:"str,omitempty"`
		Int *provider.Field `json:"int,omitempty"`
	}
	return &A{
		Str: provider.NewField("str", &a.Str),
		Int: provider.NewField("int", &a.Int),
	}
}

func (a *A) MarshalJSON() ([]byte, error) { return json.Marshal(a.toFields()) }
func (a *A) UnmarshalJSON(v []byte) error { return json.Unmarshal(v, a.toFields()) }

// B wraps the entire struct, needs to strip methods
type B A

func (b *B) toFields() interface{} {
	type stripped B // strips methods, avoids recursion
	return provider.NewField("b", (*stripped)(b))
}

func (b *B) MarshalJSON() ([]byte, error) { return json.Marshal(b.toFields()) }
func (b *B) UnmarshalJSON(v []byte) error { return json.Unmarshal(v, b.toFields()) }
