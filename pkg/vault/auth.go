package vault

import "fmt"

const (
	authTpl = "/sys/auth/%s/login"
)

type Auth interface {
	LoginK8s(role, jwt, path string) (string, error)
}

func (c *client) LoginK8s(role, jwt, path string) (string, error) {
	secret, err := c.vaultClient.Logical().Write(
		fmt.Sprintf(authTpl, path),
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
