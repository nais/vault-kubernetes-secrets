package vault

import "fmt"

const (
	root = "/sys/auth/"
)

type Auth interface {
	LoginK8s(role, jwt, path string) (string, error)
}

func (c *client) LoginK8s(role, jwt, path string) (string, error) {
	secret, err := c.vaultClient.Logical().Write(
		fmt.Sprint(root, path),
		map[string]interface{}{
			"role": role,
			"jwt":  jwt,
		},
	)
	if err != nil {
		return "", fmt.Errorf("Unable to login to vault. error: %s ", err.Error())
	}
	return secret.Auth.ClientToken, nil

}
