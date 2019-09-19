package provider_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/klaidliadon/provider"
)

func TestUnknown(t *testing.T) {
	type S struct {
		Field string
	}
	v := provider.Value{Value: &S{}}
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

type url struct {
	Str string `json:"string,omitempty"`
	Int int    `json:"int,omitempty"`
}

func (u *url) toFields() interface{} {
	type stripped url // strips methods, avoids recursion
	return provider.NewField("a", (*stripped)(u))
}

func (u *url) MarshalJSON() ([]byte, error) { return json.Marshal(u.toFields()) }
func (u *url) UnmarshalJSON(v []byte) error { return json.Unmarshal(v, u.toFields()) }

func TestURL(t *testing.T) {
	a := url{Str: "banana", Int: 42}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(a)
	})
	const address = "localhost:8888"
	go func() {
		http.ListenAndServe(address, nil)
	}()
	time.Sleep(time.Millisecond * 100)
	buff := bytes.Buffer{}
	if err := json.NewEncoder(&buff).Encode(provider.Value{
		Provider: provider.URL,
		Value:    "http://" + address,
	}); err != nil {
		t.Fatal(err)
	}
	var b url
	if err := json.NewDecoder(&buff).Decode(&b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Expected:\n%+v\nGot:\n%+v", a, b)
	}
}
