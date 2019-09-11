package storage

import (
	"context"
	"time"
)

type (
	Client interface {
		NewAuthToken(ctx context.Context, projectId int, role string, ttl time.Duration) (t Token, err error)
	}

	Token struct {
		TTL    time.Duration
		Secret string
	}
)
