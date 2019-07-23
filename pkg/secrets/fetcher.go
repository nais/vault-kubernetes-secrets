package secrets

import (
	"github.com/spf13/viper"
	"github.com/nais/vault-kubernetes-secrets/pkg/vault"
	"fmt"
	"github.com/spf13/afero"
)

const (
	envSecretsDestPath = "VKS_SECRET_DEST_PATH"
	envJwtFile         = "VKS_SERVICE_ACCOUNT_TOKEN"
)

func init() {
	viper.BindEnv(envJwtFile)
	viper.SetDefault(envJwtFile, "/var/run/secrets/kubernetes.io/serviceaccount/token")
}

func New(options vault.ClientOptions) SecretFetcher {
	fs := afero.NewOsFs()
	return secretFetcher{
		auth:         vault.NewAuthClient(options),
		kv:           vault.NewKVClient(options),
		jwtRetriever: jwtFromFile(fs, viper.GetString(envJwtFile)),
		secretWriter: writeToFile(fs, viper.GetString(envSecretsDestPath)),
	}
}

type SecretFetcher interface {
	FetchSecrets(role, authPath, kvPath string) (err error, secretsFetched int)
	FetchToken(role, authPath string) error
}

type secretFetcher struct {
	auth         vault.Auth
	kv           vault.KV
	jwtRetriever func() (string, error)
	secretWriter func(string, map[string]string) error
}

func (s secretFetcher) FetchToken(role, authPath string) error {
	tokenError, accessToken := s.fetchToken(role, authPath)
	if tokenError != nil {
		return tokenError
	}
	return s.secretWriter(accessToken, make(map[string]string, 0))
}

func (s secretFetcher) FetchSecrets(role, authPath, kvPath string) (err error, secretsFetched int) {

	tokenError, accessToken := s.fetchToken(role, authPath)
	if tokenError != nil {
		return tokenError, 0
	}

	secrets, fetchError := s.kv.Get(kvPath, accessToken)
	if fetchError != nil {
		return fetchError, 0
	}

	return s.secretWriter(accessToken, secrets), len(secrets)
}

func (s secretFetcher) fetchToken(role, authPath string) (err error, token string) {
	jwt, jwtError := s.jwtRetriever()
	if jwtError != nil {
		return jwtError, ""
	}

	accessToken, loginError := s.auth.LoginK8s(role, jwt, authPath)
	if loginError != nil {
		return loginError, ""
	}

	return nil, accessToken
}

func jwtFromFile(fs afero.Fs, jwtFile string) func() (token string, err error) {
	return func() (token string, err error) {
		if b, e := afero.ReadFile(fs, jwtFile); e != nil {
			return "", e
		} else {
			return string(b), nil
		}
	}
}

func writeToFile(fs afero.Fs, destDir string) func(token string, secrets map[string]string) error {
	return func(token string, secrets map[string]string) error {
		tokenDest := destDir + "/vault_token"
		err := afero.WriteFile(fs, tokenDest, []byte(token), 0644)
		if (err != nil) {
			return fmt.Errorf("Failed to write vault token to %s. Error: %s", tokenDest, err.Error())
		}
		for k, v := range secrets {
			dest := destDir + "/" + k
			if err = afero.WriteFile(fs, dest, []byte(v), 0644); err != nil {
				return fmt.Errorf("Fail to write secret %s to  %s. Error: ", k, err.Error())
			}
		}
		return nil
	}
}
