package secrets

import (
	"testing"
	"github.com/nais/vault-kubernetes-secrets/test/mocks"
	"github.com/nbio/st"
	"reflect"
)

func TestFetchSecrets(t *testing.T) {

	t.Run("Happy day test", func(t *testing.T) {
		var actualSecrets map[string]string
		expectedSecrets := map[string]string{"k1": "v1", "k2": "v2"}

		auth := new(mocks.Auth)
		auth.On("LoginK8s", "role", "token", "kubernetes").Return("accessToken", nil)

		kv := new(mocks.KV)
		kv.On("Get", "kv", "accessToken").Return(expectedSecrets, nil)

		s := secretFetcher{
			secretWriter: func(strings map[string]string) error {
				actualSecrets = strings
				return nil
			},
			jwtRetriever: func() (string, error) {
				return "token", nil
			},
			auth: auth,
			kv:   kv,
		}

		e := s.FetchSecrets("role", "kubernetes", "kv")

		st.Assert(t, e, nil)
		st.Assert(t, reflect.DeepEqual(actualSecrets, expectedSecrets), true)
		auth.AssertExpectations(t)
		kv.AssertExpectations(t)

	})

}
