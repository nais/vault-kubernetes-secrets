package secrets

import (
	"testing"
	"github.com/nais/vault-kubernetes-secrets/test/mocks"
	"github.com/nbio/st"
	"github.com/spf13/afero"
)

var fs = afero.NewMemMapFs()

func newFetchSecretsStub(expectedSecrets map[string]string) SecretFetcher {
	auth := new(mocks.Auth)
	auth.On("LoginK8s", "role", "token", "kubernetes").Return("accessToken", nil)

	kv := new(mocks.KV)
	kv.On("Get", "kv", "accessToken").Return(expectedSecrets, nil)

	afero.WriteFile(fs, "/token/jwt", []byte("token"), 0644)
	return secretFetcher{
		secretWriter: writeToFile(fs, "/secrets"),
		jwtRetriever: jwtFromFile(fs, "/token/jwt"),
		auth:         auth,
		kv:           kv,
	}
}
func TestFetchSecrets(t *testing.T) {
	tests := []struct {
		testName string
		secrets  map[string]string
	}{
		{"Happy day test", map[string]string{"k1": "v1", "k2": "v2"}},
		{"Handle no secrets", map[string]string{}},
	}

	for _, test := range tests {
		s := newFetchSecretsStub(test.secrets)

		e, n := s.FetchSecrets("role", "kubernetes", "kv")

		st.Assert(t, e, nil)
		st.Assert(t, n, len(test.secrets))
		for k, v := range test.secrets {
			secret, e := afero.ReadFile(fs, "/secrets/"+k)
			st.Assert(t, e, nil)
			st.Assert(t, string(secret), v)
		}

	}
}
