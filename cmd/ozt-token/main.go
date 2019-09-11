package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/hashicorp/vault/api"

	application "gitlab.ozon.ru/arubtsov/ozt-token/internal/app"
	"gitlab.ozon.ru/arubtsov/ozt-token/internal/repository_client/gitlab_client"
	"gitlab.ozon.ru/arubtsov/ozt-token/internal/storage/vault"
)

var (
	projectId    = flag.Int("projectId", 0, "gitlab projectId")
	authToken    = flag.String("gitlabToken", "", "gitlab auth token")
	gitlabHost   = flag.String("gitlab", "gitlab.ozon.ru", "gitlab host")
	vaultAddress = flag.String("vault", "", "vault address, like `http://localhost:8200`")
	vaultToken   = flag.String("vaultToken", "", "vault auth token")
	envVal       = flag.String("variable", "", "env variable name")
	env          = flag.String("env", "", "environment type [prod,stg,dev]")
	tokenTTL     = flag.Duration("ttl", 24*time.Hour, "token ttl")
)

var availableEnv = map[string]struct{}{"prod": {}, "stg": {}, "dev": {}}

func main() {
	flag.Parse()

	if *authToken == "" {
		log.Fatal("auth token not set")
	}

	if *vaultAddress == "" {
		log.Fatal("vault address not set")
	}

	if *vaultToken == "" {
		log.Fatal("vault token not set")
	}

	if *gitlabHost == "" {
		log.Fatal("gitlab host not set")
	}

	if *envVal == "" {
		log.Fatal("variable name not set")
	}

	if *env == "" {
		log.Fatal("env not set")
	}

	if _, ok := availableEnv[*env]; !ok {
		log.Fatalf("incorrect env value: [%q], available: [prod, stg, dev]", *env)
	}

	if *projectId == 0 {
		log.Fatal("projectId not set")
	}

	gitlabCli := gitlab_client.New(*gitlabHost, *authToken)

	vaultCfg := api.DefaultConfig()
	vaultCfg.Address = *vaultAddress
	vc, err := api.NewClient(vaultCfg)

	if err != nil {
		log.Fatalf("cant create vault client: [%s]", err)
	}
	vc.SetToken(*vaultToken)
	vaultClient := vault.New(vc)

	app := application.New(vaultClient, gitlabCli, *env, *tokenTTL)

	ctx := context.Background()
	err = app.StoreTokenVariable(ctx, *projectId, *envVal)
	if err != nil {
		log.Fatalf("cant store token: [%s]", err)
	} else {
		log.Println("success")
	}
}
