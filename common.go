package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var providers = map[string]Provider{}

// List of providers
const (
	Default = ""
	Static  = "static"
	URL     = "url"
)

func init() {
	Register(Default, static{})
	Register(Static, static{})
	Register(URL, url{})
}

// Register adds a provider to the list
func Register(name string, p Provider) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("duplicate provider: %s", name))
	}
	providers[name] = p
}

// Provider sets the Valueuration
type Provider interface {
	Set(data []byte, v interface{}) error
}

// static uses the data received as source
type static struct{}

// Set uses json.Unmarshal
func (s static) Set(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type url struct{}

func (u url) Set(data []byte, v interface{}) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	resp, err := http.Get(s)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
