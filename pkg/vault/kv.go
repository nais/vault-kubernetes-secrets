package vault

import (
	"fmt"
	"reflect"
)

type KV interface {
	Get(path, token string) (map[string]string, error)
}

func (c *client) Get(path, token string) (map[string]string, error) {
	c.vaultClient.SetToken(token)
	secret, e := c.vaultClient.Logical().Read(path)
	if e != nil {
		return nil, e
	}

	secrets := make(map[string]string)

	// If the map only contains two keys; "data" and "metadata", we assume it's a
	// kv store version 2.0. So return everything in the "data" key.
	if secret == nil {
		return secrets, nil
	}

	// Heuristic; check if we only get two keys ("data" and "metadata") and guess that it's kv version 2.
	if metadata, ok := secret.Data["metadata"]; ok {
		if data, ok := secret.Data["data"]; ok {
			if len(secret.Data) == 2 {
				if reflect.TypeOf(metadata).Kind() == reflect.Map && reflect.TypeOf(data).Kind() == reflect.Map {
					hello := data.(map[string]interface{})
					for k, v := range hello {
						switch val := v.(type) {
						case string:
							secrets[k] = val
						default:
							return nil, fmt.Errorf("vault: %s[%s] has type %T", path, k, val)
						}
					}
					return secrets, nil
				}
			}
		}
	}

	for k, v := range secret.Data {
		switch val := v.(type) {
		case string:
			secrets[k] = val
		default:
			return nil, fmt.Errorf("vault: %s[%s] has type %T", path, k, val)
		}
	}
	return secrets, nil
}
