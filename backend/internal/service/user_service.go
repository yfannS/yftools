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

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest, clientIP string) (string, error)
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

func (s *userService) Login(req model.LoginRequest, clientIP string) (string, error) {
	usernameKey := s.loginFailureUsernameKey(req.Username)
	usernameIPKey := s.loginFailureUsernameIPKey(req.Username, clientIP)

	if err := s.ensureLoginAllowed(usernameKey, usernameIPKey); err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		s.recordLoginFailure(usernameKey, usernameIPKey)
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.recordLoginFailure(usernameKey, usernameIPKey)
		return "", ErrInvalidCredentials
	}

	s.clearLoginFailures(usernameKey, usernameIPKey)

	token, err := appJwt.GenerateToken(user.ID, user.Username, s.jwtExpire)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	// 将登录会话存入 Redis
	if appRedis.IsAvailable() {
		expire, _ := time.ParseDuration(s.jwtExpire)
		if expire <= 0 {
			expire = 7 * 24 * time.Hour
		}
		if err := session.Set(token, &session.Session{
			UserID:   user.ID,
			Username: user.Username,
		}, expire); err != nil {
			logger.Warn("[Login] Failed to store session in Redis: %v", err)
		}
	}

	return token, nil
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

	if err := s.checkFailureThreshold(usernameKey, s.rateLimit.Login.FailUsernamePer15m); err != nil {
		return err
	}

	if err := s.checkFailureThreshold(usernameIPKey, s.rateLimit.Login.FailUsernameIPPer5m); err != nil {
		return err
	}

	return nil
}

func (s *userService) checkFailureThreshold(key string, rule config.LimitRuleConfig) error {
	if key == "" || !rule.IsEnabled() {
		return nil
	}

	state, err := ratelimit.Get(key)
	if err != nil {
		logger.Warn("[LoginRateLimit] failed to read counter key=%s: %v", key, err)
		return nil
	}

	if state.Count >= int64(rule.Limit) {
		return ErrTooManyRequests
	}

	return nil
}

func (s *userService) recordLoginFailure(usernameKey, usernameIPKey string) {
	if !s.rateLimit.Enabled {
		return
	}

	s.incrementFailureCounter(usernameKey, s.rateLimit.Login.FailUsernamePer15m)
	s.incrementFailureCounter(usernameIPKey, s.rateLimit.Login.FailUsernameIPPer5m)
}

func (s *userService) incrementFailureCounter(key string, rule config.LimitRuleConfig) {
	if key == "" || !rule.IsEnabled() {
		return
	}

	state, err := ratelimit.Increment(key, rule.Duration())
	if err != nil {
		logger.Warn("[LoginRateLimit] failed to increment counter key=%s: %v", key, err)
		return
	}

	if state.Count >= int64(rule.Limit) {
		logger.Warn("[LoginRateLimit] threshold reached key=%s count=%d limit=%d", key, state.Count, rule.Limit)
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
