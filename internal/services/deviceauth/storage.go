package deviceauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func (service *Service) acquireDecisionLock(ctx context.Context, userCode string) (func(), error) {
	key := lockKeyPrefix + hash(normalizeUserCode(userCode))
	acquired, err := service.redis.GetClient().SetNX(ctx, key, "1", 5*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, ErrAlreadyHandled.WithDetail("device authorization is already being updated")
	}

	return func() {
		_ = service.redis.Delete(ctx, key)
	}, nil
}

func (service *Service) loadState(ctx context.Context, deviceHash string) (*deviceState, error) {
	payload, err := service.redis.Get(ctx, deviceKey(deviceHash))
	if err != nil {
		return nil, ErrAuthorizationExpired
	}

	var state deviceState
	if err := json.Unmarshal([]byte(payload), &state); err != nil {
		return nil, ErrInvalidCode.WithErr(err)
	}
	if state.ExpiresAt.IsZero() || !time.Now().UTC().Before(state.ExpiresAt) {
		return nil, ErrAuthorizationExpired
	}

	return &state, nil
}

func (service *Service) saveState(ctx context.Context, deviceHash string, state deviceState, ttl time.Duration) error {
	payload, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return service.redis.Set(ctx, deviceKey(deviceHash), payload, ttl)
}

func (service *Service) uniqueUserCode(ctx context.Context) (string, error) {
	for range 10 {
		code, err := randomUserCode()
		if err != nil {
			return "", err
		}

		exists, err := service.redis.Exists(ctx, userKey(code))
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", ErrInvalidCode.WithDetail("could not generate a unique user code")
}

func randomToken() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return "fuse_" + base64.RawURLEncoding.EncodeToString(value), nil
}

func randomUserCode() (string, error) {
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	value := make([]byte, 8)
	for index := range value {
		value[index] = codeAlphabet[int(random[index])%len(codeAlphabet)]
	}

	return string(value[:4]) + "-" + string(value[4:]), nil
}

func normalizeUserCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")

	if len(value) == 8 {
		return value[:4] + "-" + value[4:]
	}

	return value
}

func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func deviceKey(deviceHash string) string {
	return deviceKeyPrefix + deviceHash
}

func userKey(userCode string) string {
	return userKeyPrefix + normalizeUserCode(userCode)
}
