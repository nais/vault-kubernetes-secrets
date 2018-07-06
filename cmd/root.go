package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/nais/vault-kubernetes-secrets/pkg/secrets"
	"log"
	"os"
	"github.com/nais/vault-kubernetes-secrets/pkg/vault"
)

const (
	envVaultAddr = "VS_VAULT_ADDR"
	envAuthPath  = "VS_AUTH_PATH"
	envKvPath    = "VS_KV_PATH"
	envVaultRole = "VS_VAULT_ROLE"
)

func init() {
	FetchCmd.PersistentFlags().StringP(envVaultAddr, "u", "http://127.0.0.1:8200", "Vault address")
	FetchCmd.PersistentFlags().StringP(envAuthPath, "a", "/kubernetes", "Kubernetes auth path")
	FetchCmd.PersistentFlags().StringP(envKvPath, "m", "/secret", "Secret backend")
	FetchCmd.PersistentFlags().StringP(envVaultRole, "r", "app", "Vault role")
	viper.BindPFlag(envVaultAddr, FetchCmd.PersistentFlags().Lookup(envVaultAddr))
	viper.BindPFlag(envAuthPath, FetchCmd.PersistentFlags().Lookup(envAuthPath))
	viper.BindPFlag(envKvPath, FetchCmd.PersistentFlags().Lookup(envKvPath))
	viper.BindPFlag(envVaultRole, FetchCmd.PersistentFlags().Lookup(envVaultRole))
	viper.BindEnv(envKvPath)
	viper.BindEnv(envAuthPath)
	viper.BindEnv(envVaultAddr)
	viper.BindEnv(envVaultRole)

}

var FetchCmd = &cobra.Command{
	Use:   "vs",
	Short: "Fetch vault secrets to file",
	Long:  " Fetch vault secrets to file",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stdout, "vs", log.LstdFlags)

		fetcher := secrets.New(vault.ClientOptions{
			Server: viper.GetString(envVaultAddr),
			Logger: logger,
		})
		role, authMount, kvMount := viper.GetString(envVaultRole), viper.GetString(envAuthPath), viper.GetString(envKvPath)
		if err := fetcher.FetchSecrets(role, authMount, kvMount); err != nil {
			log.Fatalf("Unable to fetch secrets. error: %s", err.Error())
		}

	},
}
