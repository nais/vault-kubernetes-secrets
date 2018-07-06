package vault

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"github.com/nbio/st"
)

var sampleKvResponse = `{
  "auth": null,
  "data": {
    "foo": "bar",
	"bar": "foo"
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}`

func TestKV(t *testing.T) {
	defer gock.Off()

	t.Run("Calling KV backend with a token should return secrets", func(t *testing.T) {
		token := "token"

		gock.New("http://vault.foo.bar").
			Get("/v1/secret").
			MatchHeader("X-Vault-Token", token).
			Reply(200).
			JSON(sampleKvResponse)

		c := NewKVClient(ClientOptions{
			Server: "http://vault.foo.bar",
			HTTPClient: &http.Client{
				Transport: http.DefaultTransport,
			},
		})

		secrets, e := c.Get("secret", token)
		st.Assert(t, e, nil)
		st.Assert(t, len(secrets), 2)
		st.Assert(t, secrets["foo"], "bar")
		st.Assert(t, secrets["bar"], "foo")
		st.Assert(t, gock.IsDone(), true)
	})
}
