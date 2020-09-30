package auth

import (
	"github.com/devopsfaith/krakend/config"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github.com/gosha20777/krakend-cooke-auth"

// Credentials contains the pair user:pass
type Credentials struct {
	Url string
}

// ConfigGetter extracts the credentials from the extra config details
func ConfigGetter(e config.ExtraConfig) interface{} {
	cfg, ok := e[Namespace]
	if !ok {
		return nil
	}
	data, ok := cfg.(map[string]interface{})
	if !ok {
		return nil
	}

	v, ok := data["url"]
	if !ok {
		return nil
	}

	url, ok := v.(string)
	if !ok {
		return nil
	}

	return Credentials{url}
}
