package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/nais/vault-kubernetes-secrets/pkg/secrets"
	"log"
	"os"
	"github.com/nais/vault-kubernetes-secrets/pkg/vault"
	"github.com/nais/vault-kubernetes-secrets/pkg/renewer"
)

const (
	envVaultAddr = "VKS_VAULT_ADDR"
	envAuthPath  = "VKS_AUTH_PATH"
	envKvPath    = "VKS_KV_PATH"
	envVaultRole = "VKS_VAULT_ROLE"
	envIsSidecar = "VKS_IS_SIDECAR"
	envSecretsDestPath = "VKS_SECRET_DEST_PATH"
)

func init() {
	FetchCmd.PersistentFlags().StringP(envVaultAddr, "u", "http://127.0.0.1:8200", "Vault address")
	FetchCmd.PersistentFlags().StringP(envAuthPath, "a", "/kubernetes", "Kubernetes auth path")
	FetchCmd.PersistentFlags().StringP(envKvPath, "m", "/secret", "Secret backend")
	FetchCmd.PersistentFlags().StringP(envVaultRole, "r", "app", "Vault role")
	FetchCmd.PersistentFlags().StringP(envIsSidecar, "s", "false", "Whether vks runs as a sidecar or init container")
	viper.BindPFlag(envVaultAddr, FetchCmd.PersistentFlags().Lookup(envVaultAddr))
	viper.BindPFlag(envAuthPath, FetchCmd.PersistentFlags().Lookup(envAuthPath))
	viper.BindPFlag(envKvPath, FetchCmd.PersistentFlags().Lookup(envKvPath))
	viper.BindPFlag(envVaultRole, FetchCmd.PersistentFlags().Lookup(envVaultRole))
	viper.BindPFlag(envIsSidecar, FetchCmd.PersistentFlags().Lookup(envIsSidecar))
	viper.BindEnv(envKvPath)
	viper.BindEnv(envAuthPath)
	viper.BindEnv(envVaultAddr)
	viper.BindEnv(envVaultRole)
	viper.BindEnv(envIsSidecar)
	viper.BindEnv(envSecretsDestPath)
	viper.SetDefault(envSecretsDestPath, "/var/run/secrets/naisd.io/vault")
}

var FetchCmd = &cobra.Command{
	Use:   "vs",
	Short: "Fetch vault secrets to file",
	Long:  " Fetch vault secrets to file",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stdout, "vs", log.LstdFlags)

		if (viper.GetBool(envIsSidecar)) {
			runner := renewer.New(
				viper.GetString(envVaultAddr),
				viper.GetString(envSecretsDestPath + "/vault_token"),
			)
			runner.Run()
		} else {
			fetcher := secrets.New(vault.ClientOptions{
				Server: viper.GetString(envVaultAddr),
				Logger: logger,
			})
			role, authMount, kvMount := viper.GetString(envVaultRole), viper.GetString(envAuthPath), viper.GetString(envKvPath)
			if err, secretsFetched := fetcher.FetchSecrets(role, authMount, kvMount); err != nil {
				log.Fatalf("Unable to fetch secrets. error: %s", err.Error())
			} else {
				log.Printf("Fetched %d secret(s)", secretsFetched)
			}
		}
	},
}
