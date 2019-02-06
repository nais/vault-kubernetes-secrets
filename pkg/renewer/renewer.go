package renewer

import (
	"os"
	"io/ioutil"
	"log"
	"time"
	"github.com/hashicorp/vault/api"
)

func suggestedRefreshTime(ttl float64) float64 {
	if (ttl < 60) {
		return ttl / 2
	} else {
		return ttl - 30
	}
}

type Renewer interface {
	Run()
}

type renewer struct {
	vaultAddr string
	tokenPath string
}

func New(vaultAddr string, tokenPath string) Renewer {
	return renewer {
		vaultAddr: vaultAddr,
		tokenPath: tokenPath,
	}
}

func (r renewer) Run() {
	if _, err := os.Stat(r.tokenPath); err == nil {
		data, err := ioutil.ReadFile(r.tokenPath)
		if err != nil {
			log.Fatalf("Could not read the token file %s: %s", r.tokenPath, err.Error())
		}
		token := string(data)

		c, _ := api.NewClient(&api.Config{
			Address:    r.vaultAddr,
		})
		c.SetToken(token)

		tokenMeta, err := c.Auth().Token().LookupSelf()

		if err != nil {
			log.Fatalf("Could not lookup info about the token: %s", err.Error())
			os.Exit(1)
		} else {
			var ttl time.Duration

			ttl, _ = tokenMeta.TokenTTL()

			for {
				suggested := suggestedRefreshTime(ttl.Seconds())
				time.Sleep(time.Duration(suggested) * time.Second)
				newToken, err := c.Auth().Token().RenewSelf(0)
				if err != nil {
					log.Fatalf("Could not renew the Vault token: %s", err.Error())
					os.Exit(1)
				}
				ttl, _ = newToken.TokenTTL()
				log.Printf("Renewed the Vault token, with TTL %f", ttl.Seconds())
			}
		}
	} else if os.IsNotExist(err) {
		log.Fatalf("The vault token file %s does not exist", r.tokenPath)
		os.Exit(1)
	} else {
		log.Fatalf("Could not read the vault token file %s: %s", r.tokenPath, err.Error())
		os.Exit(1)
	}
}