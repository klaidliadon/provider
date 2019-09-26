package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	Register("name", nil)
}

func TestUnknown(t *testing.T) {
	type S struct {
		Field string
	}
	v := Value{Value: &S{}}
	if err := json.Unmarshal([]byte(`{
		"provider":"unknown",
		"value":{"field":"value"}
	}`), &v); err == nil {
		t.Errorf("error expected")
	}
	if err := json.Unmarshal([]byte(`{
		"provider":"static",
		"value":{"field":"value"}
	}`), &v); err != nil {
		t.Errorf("no error expected, got %s", err)
	}
	if v := v.Value.(*S).Field; v != "value" {
		t.Errorf("%s expected, got %s", "value", v)
	}
}

type U struct {
	Str string `json:"string,omitempty"`
	Int int    `json:"int,omitempty"`
}

func (u *U) toFields() interface{} {
	type stripped U // strips methods, avoids recursion
	return NewField("a", (*stripped)(u))
}

func (u *U) MarshalJSON() ([]byte, error) { return json.Marshal(u.toFields()) }
func (u *U) UnmarshalJSON(v []byte) error { return json.Unmarshal(v, u.toFields()) }

func TestURL(t *testing.T) {
	a := U{Str: "banana", Int: 42}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(a)
	})
	const address = "localhost:8888"
	go func() {
		http.ListenAndServe(address, nil)
	}()
	time.Sleep(time.Millisecond * 100)
	buff := bytes.Buffer{}
	if err := json.NewEncoder(&buff).Encode(Value{
		Provider: URL,
		Value:    "http://" + address,
	}); err != nil {
		t.Fatal(err)
	}
	var b U
	if err := json.NewDecoder(&buff).Decode(&b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Expected:\n%+v\nGot:\n%+v", a, b)
	}
}
