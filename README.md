# vault-kubernetes-secrets

Authenticate against a Vault Kubernetes backend and fetch secrets for a KV backend.
Tailormade for use as a init container to fetch and mount secrets from Vault into a pod.

## Environment variables used

| Name   | Usage | Example(default) 
|:-------|:------|:----------------
|`VKS_VAULT_ADDR` |Address to the vault api| http://127.0.0.1:8200
|`VKS_AUTH_PATH` |The path to the Kubernetes Auth mount| /kubernetes 
|`VKS_KV_PATH` |The path to the KV mount| /secret 
|`VKS_VAULT_ROLE`|The vault role to authenticate against| app 
|`VKS_SECRET_DEST_PATH`|The path to the directory where the vault secrets will be stored| /var/run/secrets/naisd.io/vault 
|`VKS_SERVICE_ACCOUNT_TOKEN`|Path to Kubernetes service account token for which to authenticate against Vault| /var/run/secrets/kubernetes.io/serviceaccount/token 
 