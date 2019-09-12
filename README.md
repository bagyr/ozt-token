# ozt-token

Tool for semi-automated Vault <-> Gitlab CI variables sync.

## Usage:

```
make build
./bin/ozt-token --env={dev|stg|prod} \
    --projectId={gitlab_project_id} \
    --gitlabToken={gitab_access_token} \
    --variable={variable_name} \
    --vaultToken={vault_access_token} \
    --vault={vault_api_endpoint} \
    --gitlab={gitlab_api_endpoint}
```

## Local Vault instance for dry run:

```
make run.vault
make down.vault
```
access token: myroottoken

endpoint: http://localhost:8200
