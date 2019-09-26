# Provider

[![GoDoc](https://godoc.org/klaidliadon.dev/provider?status.svg)](https://godoc.org/klaidliadon.dev/provider)
[![Build Status](https://travis-ci.org/klaidliadon/provider.svg?branch=master)](https://travis-ci.org/klaidliadon/provider) 
[![codecov.io](http://codecov.io/github/klaidliadon/provider/coverage.svg?branch=master)](http://codecov.io/github/klaidliadon/provider?branch=master)

Provider is a go package that allows to use different sources to unmarshal or marshal a struct.

It allows to use a structure with `Provider` and `Value` and treat the value according to the provider specified.

## List of default providers

### Static provider

Static Provider expects athe value to be the JSON itself.

```json
{
    "provider": "static",
    "value": {
        "some": "json"
    }
}
```

### URL provider

URL Provider requires the value to be an URL, will execute the HTTP request and use the response body for the unmarshal.

```json
{
    "provider": "url",
    "value": "http://example.com/content"
}
```

## Usage

It allows to use helpers in the JSON methods to use the package capabilities. 
A type can use a `Field`, or a `struct` containing several of them in the methods.
Each `Field` will use the given provider.

```go
// Field is an helper used to inject Value
type Field struct {
    Name string
	Ptr  *interface{}
}
```

Field has a name, which is used for error message, and a pointer to the interface used for marshal and unmarshal.

Let's take the following type, which returns a structure with two fields:

```go
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
```

It will use a provider for each field, decoding the following JSON:

```json
{
    "str": {
        "provider": "static",
        "value": "banana"
    },
    "int": {
        "provider": "static",
        "value" :42
    }
}
```

This other one uses just one provider for the structure itself. It requires to create a private type to strip the methods from the structure.

```go
type B A

func (b *B) toFields() interface{} {
	type stripped B // strips methods, avoids recursion
	return provider.NewField("b", (*stripped)(b))
}

func (b *B) MarshalJSON() ([]byte, error) { return json.Marshal(b.toFields()) }
func (b *B) UnmarshalJSON(v []byte) error { return json.Unmarshal(v, b.toFields()) }
```

It will result in a JSON that looks like the following:

```json
{
    "provider": "static",
    "value": {
        "str": "banana",
        "int": 42
    }
}
```