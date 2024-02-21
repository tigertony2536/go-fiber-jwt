package core

import "github.com/tigertony2536/go-login/internal/core/domain"

type CacheRepository interface {
	SaveTokenToRedis() error
	GetTokenFromRedis() (*domain.PairedToken, error)
}
