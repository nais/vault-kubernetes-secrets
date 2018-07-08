package vault

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"github.com/nbio/st"
)

var sampleResponse = `
{
  "auth": {
    "client_token": "test_token",
    "accessor": "afa306d0-be3d-c8d2-b0d7-2676e1c0d9b4",
    "policies": [
      "default"
    ],
    "metadata": {
      "role": "test",
      "service_account_name": "vault-auth",
      "service_account_namespace": "default",
      "service_account_secret_name": "vault-auth-token-pd21c",
      "service_account_uid": "aa9aa8ff-98d0-11e7-9bb7-0800276d99bf"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
`

func TestAuth(t *testing.T) {
	defer gock.Off()

	t.Run("Calling auth backend should return a token", func(t *testing.T) {
		gock.New("http://vault.foo.bar").
			Put("/v1/auth/kubernetes/login").
			Reply(200).
			JSON(sampleResponse)

		c := NewAuthClient(ClientOptions{
			Server: "http://vault.foo.bar",
			HTTPClient: &http.Client{
				Transport: http.DefaultTransport,
			},
		})

		token, e := c.LoginK8s("role", "jwt", "/kubernetes")

		st.Assert(t, e, nil)
		st.Assert(t, token, "test_token")
		st.Assert(t, gock.IsDone(), true)

	})

}
