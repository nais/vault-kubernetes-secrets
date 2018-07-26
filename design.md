| version | application | environment | namespace       | service acc | role name   | policy                | 
|:-------:|:-----------:|:-----------:|:---------------:|:-----------:|:-----------:|:---------------------:|
| old     | app1        | default     | default         | app1        | app1-default| read kv/app1/app      |  
| new     | app1        | app         | app1            | app         | app1-app    | read kv/app1/app      |  
| old     | app1        | t1          | t1              | app1        | app1-t1     | read kv/app1/t1       |  
| new     | app1        | t1          | app1            | t1          | app1-t1     | read kv/app1/t1       |  

## Scenario 1 handle only new app namespaced. 


### Possibly multiple instances in preprod and one in prod fss.
```yaml 
  app: app1 
  env: app, q1, t1 #Optional. Defaults to "app" which means one environment.
  auth:
    ldap:
      groups:
      - "0000-ga-aura"
```

For each zone:
- For each environment create a read policy with path given by the application and environment(kv/${env}/{zone}/app1/${env}).
- For each environment create a write policy to the path given by the application and environment(kv/${env}/${zone}/app1/${env}).
- For each environment create a kubernetes auth role scoped to application namespace and  "environment" service account and grant the role the appropriate read policy.

For ldap auth:
- Grant the ldap group(s) all the write policies. We must handle the scenario where a team, group or whatever owns multipli applications.

#### Vault resources

##### Policies
``` hcl
  name: "preprod_fss_app1_app_read"
  path "kv/preprod/fss/app1/app {
    capabilities = ["read"]
  }
```
``` hcl
  name: "preprod_fss_app1_q1_read"
  path "kv/preprod/fss/app1/q1 {
    capabilities = ["read"]
  }
```

``` hcl
  name: "preprod_fss_app1_app_write"
  path "kv/preprod/fss/app1/app {
    capabilities = ["create", "read", "update", "delete"] 
  }
```
``` hcl
  name: "preprod_fss_app1_q1_write"
  path "kv/preprod/fss/app1/q1 {
    capabilities = ["create", "read", "update", "delete"] 
  }
```

```hcl
......... And so on for each environment in preprod.
```

``` hcl
  name: "prod_fss_app1_app_read"
  path "kv/prod/fss/app1/app {
    capabilities = ["read"]
  }
```

``` hcl
  name: "prod_fss_app1_app_write"
  path "kv/prod/fss/app1/app {
    capabilities = ["create", "read", "update", "delete"] 
  }
```

##### Roles

```hcl
 path "auth/kubernetes/preprod/fss/role/app1-app {
   bound_service_account_names = [app]
   bound_service_account_namespaces = [app1]
   policies = [preprod_fss_app1_app_read]
 }
```

```hcl
 path "auth/kubernetes/preprod/fss/role/app1-q1 {
   bound_service_account_names = [q1]
   bound_service_account_namespaces = [app1]
   policies = [preprod_fss_app1_q1_read]
 }
```

```hcl
......... And so on for each environment in preprod.
```

```hcl
 path "auth/kubernetes/prod/fss/role/app1-app {
   bound_service_account_names = [app]
   bound_service_account_namespaces = [app1]
   policies = [prod_fss_app1_app_read]
 }
```

```hcl
 path "auth/ldap/groups/0000-ga-aura"{
   policies = ["preprod_fss_app1_app_write", "preprod_fss_app1_q1_write", ...., "prod_fss_app1_app_write"]
 }
```


