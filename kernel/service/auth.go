package service

import (
	"context"
	"errors"
	"fmt"
	"labflow/apperr"
	"labflow/auth"
	"labflow/models"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RegisterParams struct {
	Username    string
	Password    string
	Email       string
	PhoneNumber string
	Name        string
	Company     string
	Address     string
}

type RegisterResult struct {
	ID          int
	Xid         uuid.UUID
	Username    string
	Email       string
	PhoneNumber string
	Name        string
	Company     string
	Address     string
	CreatedAt   time.Time
	Message     string
}

type LoginParams struct {
	Identifier string
	Password   string
}

type LoginResult struct {
	Message string
	User    UserInfo
}

type UserInfo struct {
	Xid         uuid.UUID
	Username    string
	Email       string
	Name        string
	PhoneNumber string
	Company     string
}

type VerifyOTPParams struct {
	Email string
	Otp   string
}

type ResendOTPParams struct {
	Email string
}

type MessageResult struct {
	Message string
}

type MeResult struct {
	ID          int
	Xid         uuid.UUID
	Username    string
	Email       string
	Name        string
	PhoneNumber string
	Company     string
	Address     string
	Avatar      *string
	Roles       []string
	Permissions []string
}

type authRepo interface {
	CreateUser(ctx context.Context, username, password, email, phoneNumber, name, company, address string, roleID int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserVerifiedStatus(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	UpdateUserPassword(ctx context.Context, email, hashedPassword string) (*models.User, error)
}

type otpStore interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

type verificationMailer interface {
	SendVerificationEmail(to, otp string) error
}

type tokenManager interface {
	GenerateJWT(ctx context.Context, auth *models.Auth) (string, error)
	InvalidateJWT(ctx context.Context, token string) error
	GenerateResetToken(ctx context.Context, email string) (string, error)
	ParseResetToken(ctx context.Context, tokenString string) (string, error)
}

type Service struct {
	repo           authRepo
	redis          otpStore
	mailer         verificationMailer
	tokenUtils     tokenManager
	customerRoleID int
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidVerification = errors.New("user not verified, Please verify your email")
var ErrInvalidUser = errors.New("invalid user")
var ErrAlreadyVerified = errors.New("email already verified")
var ErrInvalidOtp = errors.New("invalid or expired OTP")
var ErrSendingEmail = errors.New("failed to send email")
var ErrRedisFailed = errors.New("failed to run redis service")
var ErrEmailTaken = errors.New("email or username already taken")
var ErrInvalidResetToken = errors.New("invalid or expired reset token")
var ErrOtpCooldown = errors.New("please wait before resending OTP")

func NewService(repo authRepo, redisClient otpStore, mailer verificationMailer, tokenUtils tokenManager, customerRoleID int) *Service {
	return &Service{
		repo:           repo,
		redis:          redisClient,
		mailer:         mailer,
		tokenUtils:     tokenUtils,
		customerRoleID: customerRoleID,
	}
}

func (s *Service) Register(ctx context.Context, p RegisterParams) (RegisterResult, error) {
	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		return RegisterResult{}, err
	}

	user, err := s.repo.CreateUser(ctx,
		strings.TrimSpace(p.Username),
		hashedPassword,
		strings.ToLower(strings.TrimSpace(p.Email)),
		p.PhoneNumber,
		strings.TrimSpace(p.Name),
		p.Company,
		p.Address,
		s.customerRoleID,
	)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicateKey) {
			return RegisterResult{}, ErrEmailTaken
		}
		return RegisterResult{}, err
	}

	otp := auth.GenerateOTP()
	if err := s.redis.Set(ctx, fmt.Sprintf("otp:verify:%s", user.Email), otp, 5*time.Minute); err != nil {
		return RegisterResult{}, ErrRedisFailed
	}
	if err := s.mailer.SendVerificationEmail(user.Email, otp); err != nil {
		return RegisterResult{}, ErrSendingEmail
	}

	return RegisterResult{
		ID:          user.ID,
		Xid:         user.XID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		Company:     user.Company,
		Address:     user.Address,
		CreatedAt: func(t *time.Time) time.Time {
			if t == nil {
				return time.Time{}
			}
			return *t
		}(user.CreatedAt),
		Message: "Verification email sent. Please check your inbox",
	}, nil
}

func (s *Service) LoginWithCookie(ctx context.Context, p LoginParams) (LoginResult, string, error) {
	var user *models.User
	var err error

	if strings.Contains(p.Identifier, "@") {
		user, err = s.repo.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(p.Identifier)))
	} else {
		user, err = s.repo.GetUserByUsername(ctx, strings.TrimSpace(p.Identifier))
	}
	if err != nil {
		return LoginResult{}, "", ErrInvalidCredentials
	}

	if err := auth.ComparePassword(user.Password, p.Password); err != nil {
		return LoginResult{}, "", ErrInvalidCredentials
	}

	if user.IsVerified == nil || !*user.IsVerified {
		return LoginResult{}, "", ErrInvalidVerification
	}

	token, err := s.tokenUtils.GenerateJWT(ctx, &models.Auth{ID: user.ID})
	if err != nil {
		return LoginResult{}, "", err
	}

	return LoginResult{
		Message: "Login success",
		User: UserInfo{
			Xid:         user.XID,
			Username:    user.Username,
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Company:     user.Company,
		},
	}, token, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if err := s.tokenUtils.InvalidateJWT(ctx, token); err != nil {
		return err
	}
	return nil
}

func (s *Service) VerifyOtp(ctx context.Context, p VerifyOTPParams) (MessageResult, error) {
	storedOtp, err := s.redis.Get(ctx, fmt.Sprintf("otp:verify:%s", p.Email))
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return MessageResult{}, ErrInvalidOtp
		}
		return MessageResult{}, err
	}

	if storedOtp != p.Otp {
		return MessageResult{}, ErrInvalidOtp
	}

	_, err = s.repo.UpdateUserVerifiedStatus(ctx, p.Email)
	if err != nil {
		return MessageResult{}, err
	}

	if err := s.redis.Del(ctx, fmt.Sprintf("otp:verify:%s", p.Email)); err != nil {
		return MessageResult{}, err
	}

	return MessageResult{
		Message: "Email verification success!",
	}, nil
}

func (s *Service) ResendOtp(ctx context.Context, p ResendOTPParams) (MessageResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(p.Email)))
	if err != nil {
		return MessageResult{}, ErrInvalidUser
	}

	if user.IsVerified != nil && *user.IsVerified {
		return MessageResult{}, ErrAlreadyVerified
	}

	key := fmt.Sprintf("otp:verify:%s", user.Email)
	ttl, err := s.redis.TTL(ctx, key)
	if err == nil && ttl > 4*time.Minute {
		return MessageResult{}, ErrOtpCooldown
	}

	otp := auth.GenerateOTP()
	if err := s.redis.Set(ctx, fmt.Sprintf("otp:verify:%s", p.Email), otp, 5*time.Minute); err != nil {
		return MessageResult{}, err
	}

	if err := s.mailer.SendVerificationEmail(p.Email, otp); err != nil {
		log.Printf("Error: %v", err)
		return MessageResult{}, ErrSendingEmail
	}

	return MessageResult{
		Message: "OTP sent!",
	}, nil
}

func (s *Service) GetMe(ctx context.Context, userID int) (MeResult, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return MeResult{}, err
	}

	roles := make([]string, 0, len(user.Roles))
	permSet := make(map[string]struct{})
	for _, role := range user.Roles {
		if role.RoleName != nil {
			roles = append(roles, *role.RoleName)
		}
		for _, access := range role.Accesses {
			if access.AccessName != nil {
				permSet[*access.AccessName] = struct{}{}
			}
		}
	}
	permissions := make([]string, 0, len(permSet))
	for p := range permSet {
		permissions = append(permissions, p)
	}

	return MeResult{
		ID:          user.ID,
		Xid:         user.XID,
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Company:     user.Company,
		Address:     user.Address,
		Avatar:      FileURL(FileCategoryUser, user.Avatar),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return ErrInvalidUser
	}

	otp := auth.GenerateOTP()
	if err := s.redis.Set(ctx, fmt.Sprintf("otp:reset:%s", user.Email), otp, 5*time.Minute); err != nil {
		return err
	}

	if err := s.mailer.SendVerificationEmail(user.Email, otp); err != nil {
		log.Printf("Error sending forgot password OTP email: %v", err)
		return ErrSendingEmail
	}

	return nil
}

func (s *Service) VerifyResetOtp(ctx context.Context, email, otp string) (string, error) {
	emailNormalized := strings.ToLower(strings.TrimSpace(email))
	storedOtp, err := s.redis.Get(ctx, fmt.Sprintf("otp:reset:%s", emailNormalized))
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return "", ErrInvalidOtp
		}
		return "", err
	}

	if storedOtp != otp {
		return "", ErrInvalidOtp
	}

	if err := s.redis.Del(ctx, fmt.Sprintf("otp:reset:%s", emailNormalized)); err != nil {
		return "", err
	}

	token, err := s.tokenUtils.GenerateResetToken(ctx, emailNormalized)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) (MessageResult, error) {
	email, err := s.tokenUtils.ParseResetToken(ctx, token)
	if err != nil {
		return MessageResult{}, ErrInvalidResetToken
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return MessageResult{}, fmt.Errorf("failed to secure password: %w", err)
	}

	_, err = s.repo.UpdateUserPassword(ctx, email, hashedPassword)
	if err != nil {
		return MessageResult{}, err
	}

	return MessageResult{
		Message: "Password reset successful. Please login again with your new password.",
	}, nil
}
