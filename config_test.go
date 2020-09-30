package auth

import (
	"testing"

	"github.com/devopsfaith/krakend/config"
)

func TestConfigGetter(t *testing.T) {
	v := ConfigGetter(config.ExtraConfig(map[string]interface{}{Namespace: map[string]interface{}{"url": "a", "cookie": "b"}}))
	if v == nil {
		t.Fail()
	}
	credentials, ok := v.(Credentials)
	if !ok {
		t.Fail()
	}
	if credentials.Url != "a" {
		t.Fail()
	}
	if credentials.Cookie != "b" {
		t.Fail()
	}

	if v = ConfigGetter(config.ExtraConfig(map[string]interface{}{Namespace: map[string]interface{}{}})); v != nil {
		t.Fail()
	}

	if v = ConfigGetter(config.ExtraConfig(map[string]interface{}{Namespace: map[string]interface{}{"url": true}})); v != nil {
		t.Fail()
	}

	if v = ConfigGetter(config.ExtraConfig(map[string]interface{}{Namespace: map[string]interface{}{"url": 42}})); v != nil {
		t.Fail()
	}

	if v = ConfigGetter(config.ExtraConfig(map[string]interface{}{Namespace: true})); v != nil {
		t.Fail()
	}

	if v = ConfigGetter(config.ExtraConfig(map[string]interface{}{"url": "a"})); v != nil {
		t.Fail()
	}
}
