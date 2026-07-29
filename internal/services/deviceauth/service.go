package deviceauth

import (
	"context"
	"strings"
	"time"

	domainCLI "fuse/internal/domain/cli"
	"fuse/internal/domain/user"
	redisInfrastructure "fuse/internal/infrastructure/redis"
	"fuse/pkg/config"
)

const (
	deviceTTL       = 10 * time.Minute
	pollInterval    = 5
	deviceKeyPrefix = "auth:device:code:"
	userKeyPrefix   = "auth:device:user:"
	rateKeyPrefix   = "auth:device:rate:"
	lockKeyPrefix   = "auth:device:lock:"
)

type Service struct {
	redis       *redisInfrastructure.RedisClient
	credentials domainCLI.Repository
	users       user.Repository
	frontendURL string
}

func NewService(cfg *config.Config, redis *redisInfrastructure.RedisClient, credentials domainCLI.Repository, users user.Repository) *Service {
	return &Service{
		redis:       redis,
		credentials: credentials,
		users:       users,
		frontendURL: strings.TrimRight(cfg.Frontend.URL, "/"),
	}
}

func (service *Service) Allow(ctx context.Context, key string, limit int64, window time.Duration) bool {
	redisKey := rateKeyPrefix + hash(key)
	count, err := service.redis.GetClient().Incr(ctx, redisKey).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		_ = service.redis.GetClient().Expire(ctx, redisKey, window).Err()
	}

	return count <= limit
}
