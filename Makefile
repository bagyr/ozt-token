run.vault:
	 docker run --rm -d --name local-vault --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG={"backend": {"file": {"path": "/vault/file"}}, "default_lease_ttl": "168h", "max_lease_ttl": "720h"}' -e 'VAULT_DEV_ROOT_TOKEN_ID=myroottoken' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' -p 8200:8200 vault

down.vault:
	docker stop local-vault

build:
	go build -o ./bin/ozt-token cmd/ozt-token/main.go