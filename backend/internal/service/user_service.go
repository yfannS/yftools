package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"md2html/internal/config"
	"md2html/internal/model"
	"md2html/internal/repository"
	appJwt "md2html/pkg/jwt"
	"md2html/pkg/logger"
	"md2html/pkg/ratelimit"
	appRedis "md2html/pkg/redis"
	"md2html/pkg/session"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var ErrTooManyRequests = errors.New("请求过于频繁，请稍后再试")
var ErrUsernameExists = errors.New("用户名已存在")
var ErrInvalidCredentials = errors.New("用户名或密码错误")

type LoginFailureError struct {
	message string
	cause   error
	data    model.LoginFailureData
}

func (e *LoginFailureError) Error() string {
	return e.message
}

func (e *LoginFailureError) Data() model.LoginFailureData {
	return e.data
}

func (e *LoginFailureError) Unwrap() error {
	return e.cause
}

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest, clientIP string) (*model.LoginResponse, error)
	GetProfile(userID int64) (*model.UserProfile, error)
	Logout(token string) error
}

type userService struct {
	userRepo  repository.UserRepository
	jwtExpire string
	rateLimit config.RateLimitConfig
}

func NewUserService(userRepo repository.UserRepository, jwtExpire string, rateLimit config.RateLimitConfig) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtExpire: jwtExpire,
		rateLimit: rateLimit,
	}
}

func (s *userService) Register(req model.RegisterRequest) error {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return ErrUsernameExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = s.userRepo.Create(req.Username, string(hash))
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrUsernameExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *userService) Login(req model.LoginRequest, clientIP string) (*model.LoginResponse, error) {
	usernameKey := s.loginFailureUsernameKey(req.Username)
	usernameIPKey := s.loginFailureUsernameIPKey(req.Username, clientIP)

	if err := s.ensureLoginAllowed(usernameKey, usernameIPKey); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, s.recordLoginFailure(usernameKey, usernameIPKey)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, s.recordLoginFailure(usernameKey, usernameIPKey)
	}

	s.clearLoginFailures(usernameKey, usernameIPKey)

	token, err := appJwt.GenerateToken(user.ID, user.Username, s.jwtExpire)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	expire := appJwt.ResolveExpire(s.jwtExpire)
	expiresAt := time.Now().Add(expire)

	// 将登录会话存入 Redis
	if appRedis.IsAvailable() {
		if err := session.Set(token, &session.Session{
			UserID:   user.ID,
			Username: user.Username,
		}, expire); err != nil {
			logger.Warn("[Login] Failed to store session in Redis: %v", err)
		}
	}

	return &model.LoginResponse{
		Token:     token,
		Username:  user.Username,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *userService) GetProfile(userID int64) (*model.UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return &model.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) Logout(token string) error {
	if appRedis.IsAvailable() {
		if err := session.Delete(token); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
	}
	return nil
}

func (s *userService) ensureLoginAllowed(usernameKey, usernameIPKey string) error {
	if !s.rateLimit.Enabled {
		return nil
	}

	if blocked, data := s.checkFailureThreshold(usernameKey, s.rateLimit.Login.FailUsernamePer15m); blocked {
		return &LoginFailureError{message: ErrTooManyRequests.Error(), cause: ErrTooManyRequests, data: data}
	}

	if blocked, data := s.checkFailureThreshold(usernameIPKey, s.rateLimit.Login.FailUsernameIPPer5m); blocked {
		return &LoginFailureError{message: ErrTooManyRequests.Error(), cause: ErrTooManyRequests, data: data}
	}

	return nil
}

func (s *userService) checkFailureThreshold(key string, rule config.LimitRuleConfig) (bool, model.LoginFailureData) {
	if key == "" || !rule.IsEnabled() {
		return false, model.LoginFailureData{}
	}

	state, err := ratelimit.Get(key)
	if err != nil {
		logger.Warn("[LoginRateLimit] failed to read counter key=%s: %v", key, err)
		return false, model.LoginFailureData{}
	}

	if state.Count >= int64(rule.Limit) {
		return true, model.LoginFailureData{
			RemainingAttempts: 0,
			RetryAfterSeconds: ttlSeconds(state.TTL),
			Locked:            true,
		}
	}

	return false, model.LoginFailureData{}
}

func (s *userService) recordLoginFailure(usernameKey, usernameIPKey string) error {
	if !s.rateLimit.Enabled {
		return &LoginFailureError{message: ErrInvalidCredentials.Error(), cause: ErrInvalidCredentials, data: model.LoginFailureData{}}
	}

	usernameData := s.incrementFailureCounter(usernameKey, s.rateLimit.Login.FailUsernamePer15m)
	usernameIPData := s.incrementFailureCounter(usernameIPKey, s.rateLimit.Login.FailUsernameIPPer5m)
	data := mergeLoginFailureData(usernameData, usernameIPData)

	if data.Locked {
		return &LoginFailureError{message: ErrTooManyRequests.Error(), cause: ErrTooManyRequests, data: data}
	}

	return &LoginFailureError{message: ErrInvalidCredentials.Error(), cause: ErrInvalidCredentials, data: data}
}

func (s *userService) incrementFailureCounter(key string, rule config.LimitRuleConfig) model.LoginFailureData {
	if key == "" || !rule.IsEnabled() {
		return model.LoginFailureData{RemainingAttempts: -1}
	}

	state, err := ratelimit.Increment(key, rule.Duration())
	if err != nil {
		logger.Warn("[LoginRateLimit] failed to increment counter key=%s: %v", key, err)
		return model.LoginFailureData{RemainingAttempts: -1}
	}

	remaining := rule.Limit - int(state.Count)
	if remaining < 0 {
		remaining = 0
	}

	if state.Count >= int64(rule.Limit) {
		logger.Warn("[LoginRateLimit] threshold reached key=%s count=%d limit=%d", key, state.Count, rule.Limit)
		return model.LoginFailureData{
			RemainingAttempts: 0,
			RetryAfterSeconds: ttlSeconds(state.TTL),
			Locked:            true,
		}
	}

	return model.LoginFailureData{
		RemainingAttempts: remaining,
	}
}

func (s *userService) clearLoginFailures(usernameKey, usernameIPKey string) {
	if !s.rateLimit.Enabled {
		return
	}

	keys := make([]string, 0, 2)
	if usernameKey != "" {
		keys = append(keys, usernameKey)
	}
	if usernameIPKey != "" {
		keys = append(keys, usernameIPKey)
	}
	if len(keys) == 0 {
		return
	}

	if err := ratelimit.Reset(keys...); err != nil {
		logger.Warn("[LoginRateLimit] failed to clear counters: %v", err)
	}
}

func (s *userService) loginFailureUsernameKey(username string) string {
	if !s.rateLimit.Enabled || !s.rateLimit.Login.FailUsernamePer15m.IsEnabled() {
		return ""
	}
	return fmt.Sprintf("%s:login_fail:username:%s", s.rateLimit.KeyPrefix(), normalizeRateLimitKey(username))
}

func (s *userService) loginFailureUsernameIPKey(username, clientIP string) string {
	if !s.rateLimit.Enabled || !s.rateLimit.Login.FailUsernameIPPer5m.IsEnabled() {
		return ""
	}
	if clientIP == "" {
		return ""
	}
	return fmt.Sprintf("%s:login_fail:username_ip:%s:%s", s.rateLimit.KeyPrefix(), normalizeRateLimitKey(username), normalizeRateLimitKey(clientIP))
}

func normalizeRateLimitKey(value string) string {
	return url.QueryEscape(strings.ToLower(strings.TrimSpace(value)))
}

func mergeLoginFailureData(items ...model.LoginFailureData) model.LoginFailureData {
	result := model.LoginFailureData{
		RemainingAttempts: -1,
	}

	for _, item := range items {
		if item.RemainingAttempts >= 0 && (result.RemainingAttempts == -1 || item.RemainingAttempts < result.RemainingAttempts) {
			result.RemainingAttempts = item.RemainingAttempts
		}
		if item.RetryAfterSeconds > result.RetryAfterSeconds {
			result.RetryAfterSeconds = item.RetryAfterSeconds
		}
		if item.Locked {
			result.Locked = true
		}
	}

	if result.RemainingAttempts < 0 {
		result.RemainingAttempts = 0
	}

	return result
}

func ttlSeconds(ttl time.Duration) int64 {
	if ttl <= 0 {
		return 0
	}

	seconds := int64(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	return seconds
}
