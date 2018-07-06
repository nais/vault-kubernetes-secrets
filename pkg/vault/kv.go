package vault

import (
	"fmt"
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
