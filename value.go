package provider

import (
	"encoding/json"
	"fmt"
)

// New returns a new Value
func New(v interface{}) Value {
	return Value{Provider: Default, Value: v}
}

// Value is wrapper that uses registered providers in Unmarshal
type Value struct {
	Provider string      `json:"provider,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// UnmarshalJSON uses the Provider to set the value
func (v *Value) UnmarshalJSON(data []byte) error {
	var a struct {
		Provider string          `json:"provider,omitempty"`
		Value    json.RawMessage `json:"value,omitempty"`
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	p, ok := providers[a.Provider]
	if !ok {
		return fmt.Errorf("unknown provider: %s", a.Provider)
	}
	if err := p.Set(a.Value, &v.Value); err != nil {
		return err
	}
	v.Provider = a.Provider
	return nil
}

// NewField returns a new Field
func NewField(name string, ptr interface{}) *Field {
	return &Field{
		Name: name,
		Ptr:  &ptr,
	}
}

// Field is an helper used to inject Value
type Field struct {
	Name string
	Ptr  *interface{}
}

// MarshalJSON injects the Value
func (f *Field) MarshalJSON() ([]byte, error) {
	if *f.Ptr == nil {
		return nil, nil
	}
	b, err := json.Marshal(&Value{Provider: Static, Value: *f.Ptr})
	if err != nil {
		return nil, fmt.Errorf("%s: %s\n%+v", f.Name, err, *f.Ptr)
	}
	return b, nil
}

// UnmarshalJSON uses Value to populate Fields
func (f *Field) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	cfg := New(*f.Ptr)
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("%s: %s\n%s", f.Name, err, data)
	}
	*f.Ptr = cfg.Value
	return nil
}
