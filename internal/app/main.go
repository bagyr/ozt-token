package app

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.ru/arubtsov/ozt-token/internal/repository_client"
	"gitlab.ozon.ru/arubtsov/ozt-token/internal/storage"
)

const roleTpl = "ro-%s-travel-flight-gds-pyton-api-4108"

type App struct {
	vault    storage.Client
	reposCli repository_client.Client
	tokenTTL time.Duration
	env      string
}

func New(v storage.Client, c repository_client.Client, env string, ttl time.Duration) *App {
	return &App{vault: v, reposCli: c, tokenTTL: ttl, env: env}
}

func (app *App) StoreTokenVariable(ctx context.Context, projectId int, envVar string) error {
	role := fmt.Sprintf(roleTpl, app.env)
	t, err := app.vault.NewAuthToken(ctx, projectId, role, app.tokenTTL)
	if err != nil {
		return err
	}

	vars, err := app.reposCli.Variables(ctx, projectId)
	if err != nil {
		return err
	}

	for _, v := range vars {
		if v.Key == envVar {
			v.Val = t.Secret
			return app.reposCli.UpdateVariable(ctx, projectId, v)
		}
	}

	v := repository_client.NewVar(envVar, t.Secret, true)
	return app.reposCli.CreateVariable(ctx, projectId, v)
}
