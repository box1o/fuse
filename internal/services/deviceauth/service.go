package deviceauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fuse/internal/domain/user"
	redispkg "fuse/internal/infrastructure/redis"
	computeService "fuse/internal/services/compute"
	"fuse/pkg/config"
	"fuse/pkg/deviceauthapi"
	appErrors "fuse/pkg/errors"

	"github.com/google/uuid"
)

const (
	deviceTTL    = 10 * time.Minute
	pollInterval = 5
)

var (
	ErrInvalidCode          = appErrors.New("INVALID_DEVICE_CODE", "device code is invalid")
	ErrAuthorizationPending = appErrors.New("DEVICE_AUTHORIZATION_PENDING", "authorization is pending")
	ErrAuthorizationDenied  = appErrors.New("DEVICE_AUTHORIZATION_DENIED", "authorization was denied")
	ErrAuthorizationExpired = appErrors.New("DEVICE_AUTHORIZATION_EXPIRED", "device authorization expired")
)

type State struct {
	UserCode   string    `json:"user_code"`
	ClientName string    `json:"client_name"`
	Status     string    `json:"status"`
	OwnerID    string    `json:"owner_id,omitempty"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type CodeResponse = deviceauthapi.CodeResponse
type TokenResponse = deviceauthapi.TokenResponse

type Service struct {
	redis       *redispkg.RedisClient
	compute     *computeService.Service
	users       user.Repository
	frontendURL string
}

func NewService(cfg *config.Config, redis *redispkg.RedisClient, compute *computeService.Service, users user.Repository) *Service {
	return &Service{redis: redis, compute: compute, users: users, frontendURL: strings.TrimRight(cfg.Frontend.URL, "/")}
}

func (s *Service) Create(ctx context.Context, clientName string) (*CodeResponse, error) {
	deviceCode, err := randomSecret(32)
	if err != nil {
		return nil, err
	}
	userCode, err := randomUserCode()
	if err != nil {
		return nil, err
	}

	state := State{UserCode: userCode, ClientName: strings.TrimSpace(clientName), Status: "pending", ExpiresAt: time.Now().UTC().Add(deviceTTL)}
	if state.ClientName == "" {
		state.ClientName = "Fuse CLI"
	}
	payload, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	if err := s.redis.Set(ctx, deviceKey(deviceCode), payload, deviceTTL); err != nil {
		return nil, err
	}
	if err := s.redis.Set(ctx, userKey(userCode), deviceCode, deviceTTL); err != nil {
		_ = s.redis.Delete(ctx, deviceKey(deviceCode))
		return nil, err
	}

	verificationURI := s.frontendURL + "/device"
	return &CodeResponse{
		DeviceCode: deviceCode, UserCode: userCode, VerificationURI: verificationURI,
		VerificationURIComplete: verificationURI + "?code=" + userCode,
		ExpiresIn:               int(deviceTTL.Seconds()), Interval: pollInterval,
	}, nil
}

func (s *Service) Approve(ctx context.Context, userCode string, ownerID uuid.UUID) error {
	return s.updateState(ctx, userCode, "approved", ownerID)
}

func (s *Service) Deny(ctx context.Context, userCode string, ownerID uuid.UUID) error {
	return s.updateState(ctx, userCode, "denied", ownerID)
}

func (s *Service) Exchange(ctx context.Context, deviceCode string) (*TokenResponse, error) {
	state, err := s.getState(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().After(state.ExpiresAt) {
		return nil, ErrAuthorizationExpired
	}
	switch state.Status {
	case "pending":
		return nil, ErrAuthorizationPending
	case "denied":
		return nil, ErrAuthorizationDenied
	case "approved":
	default:
		return nil, ErrInvalidCode
	}

	ownerID, err := uuid.Parse(state.OwnerID)
	if err != nil {
		return nil, ErrInvalidCode
	}
	owner, err := s.users.FindByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	token, credential, err := s.compute.IssueCLICredential(ctx, ownerID, state.ClientName)
	if err != nil {
		return nil, err
	}
	_ = s.redis.Delete(ctx, deviceKey(deviceCode), userKey(state.UserCode))
	return &TokenResponse{
		AccessToken: token, TokenType: "Bearer", ExpiresAt: credential.ExpiresAt,
		OwnerID: ownerID.String(), OwnerName: owner.Name, OwnerEmail: owner.Email,
	}, nil
}

func (s *Service) Inspect(ctx context.Context, userCode string) (*State, error) {
	deviceCode, err := s.redis.Get(ctx, userKey(normalizeUserCode(userCode)))
	if err != nil {
		return nil, ErrAuthorizationExpired
	}
	return s.getState(ctx, deviceCode)
}

func (s *Service) Owner(ctx context.Context, ownerID uuid.UUID) (*user.User, error) {
	return s.users.FindByID(ctx, ownerID)
}

func (s *Service) Allow(ctx context.Context, key string, limit int64, window time.Duration) bool {
	redisKey := "compute:device:rate:" + key
	count, err := s.redis.GetClient().Incr(ctx, redisKey).Result()
	if err != nil {
		return false
	}
	if count == 1 {
		_ = s.redis.GetClient().Expire(ctx, redisKey, window).Err()
	}
	return count <= limit
}

func (s *Service) updateState(ctx context.Context, userCode, status string, ownerID uuid.UUID) error {
	if ownerID == uuid.Nil {
		return ErrInvalidCode
	}
	userCode = normalizeUserCode(userCode)
	deviceCode, err := s.redis.Get(ctx, userKey(userCode))
	if err != nil {
		return ErrAuthorizationExpired
	}
	state, err := s.getState(ctx, deviceCode)
	if err != nil {
		return err
	}
	if state.Status != "pending" {
		return ErrInvalidCode.WithDetail("device request has already been handled")
	}
	state.Status = status
	state.OwnerID = ownerID.String()
	payload, err := json.Marshal(state)
	if err != nil {
		return err
	}
	remaining := time.Until(state.ExpiresAt)
	if remaining <= 0 {
		return ErrAuthorizationExpired
	}
	return s.redis.Set(ctx, deviceKey(deviceCode), payload, remaining)
}

func (s *Service) getState(ctx context.Context, deviceCode string) (*State, error) {
	payload, err := s.redis.Get(ctx, deviceKey(deviceCode))
	if err != nil {
		return nil, ErrAuthorizationExpired
	}
	var state State
	if err := json.Unmarshal([]byte(payload), &state); err != nil {
		return nil, fmt.Errorf("decode device authorization state: %w", err)
	}
	return &state, nil
}

func randomSecret(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func randomUserCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	value := make([]byte, 8)
	random := make([]byte, len(value))
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	for i := range value {
		value[i] = alphabet[int(random[i])%len(alphabet)]
	}
	return string(value[:4]) + "-" + string(value[4:]), nil
}

func normalizeUserCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "")
	if len(value) == 8 {
		return value[:4] + "-" + value[4:]
	}
	return value
}

func deviceKey(code string) string { return "compute:device:code:" + code }
func userKey(code string) string   { return "compute:device:user:" + normalizeUserCode(code) }
