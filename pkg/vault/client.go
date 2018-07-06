package vault

import (
	"time"
	"log"
	"io/ioutil"
	"github.com/hashicorp/vault/api"
	"net/http"
)

type ClientOptions struct {
	Server      string
	HTTPClient  *http.Client
	HTTPTimeout time.Duration
	Logger      *log.Logger
}

type Client interface {
	Auth
	KV
}

func newClient(opts ClientOptions) Client {
	c, _ := api.NewClient(&api.Config{
		Address:    opts.Server,
		Timeout:    opts.HTTPTimeout,
		HttpClient: opts.HTTPClient,
	})

	if opts.Logger == nil {
		opts.Logger = log.New(ioutil.Discard, "", 0)
	}

	return &client{
		opts:        opts,
		vaultClient: c,
	}
}
func NewAuthClient(options ClientOptions) Auth {
	return newClient(options)
}

func NewKVClient(options ClientOptions) KV {
	return newClient(options)
}

type client struct {
	opts        ClientOptions
	vaultClient *api.Client
}
