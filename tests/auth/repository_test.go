package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-auth/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	ts      *tests.TestSuite
	repo    repository.Repositories
	gormDB  *gorm.DB
	userID  int
	email   string
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := ts.GormDB()
	s.Require().NoError(err)
	s.gormDB = gormDB

	log, _ := logger.NewLogger("test", nil)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(ts.RedisClient(), log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	depsClients := tests.StartDependencyClients(gormDB, cacheStore, log, obs)
	s.repo = repository.NewRepositories(gormDB, depsClients.UserQuery, depsClients.UserCommand, depsClients.RoleQuery, depsClients.RoleCommand, depsClients.UserRole)
	s.email = "auth.repo.test@example.com"
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	tests.StopDependencyClients()
	s.ts.Teardown()
}

func (s *AuthRepositoryTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Repo",
		Email:           s.email,
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "123456",
		IsVerified:      false,
	}

	res, err := s.repo.User.CreateUser(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(s.email, res.Email)
	s.userID = int(res.UserID)
}

func (s *AuthRepositoryTestSuite) Test2_FindByEmail() {
	s.Require().NotEmpty(s.email)
	ctx := context.Background()

	found, err := s.repo.User.FindByEmail(ctx, s.email)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.userID), found.UserID)
}

func (s *AuthRepositoryTestSuite) Test3_FindById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.User.FindById(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.userID), found.UserID)
}

func (s *AuthRepositoryTestSuite) Test4_UpdateVerification() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updated, err := s.repo.User.UpdateUserIsVerified(ctx, s.userID, true)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(s.userID), updated.UserID)
}

func (s *AuthRepositoryTestSuite) Test5_UpdatePassword() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updated, err := s.repo.User.UpdateUserPassword(ctx, s.userID, "newpassword123")
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(s.userID), updated.UserID)

	// The user gRPC response never carries the password, so confirm the write
	// landed by reading the users table directly.
	var stored string
	err = s.gormDB.WithContext(ctx).
		Raw("SELECT password FROM users WHERE user_id = ?", s.userID).
		Scan(&stored).Error
	s.Require().NoError(err)
	s.Equal("newpassword123", stored)
}

func (s *AuthRepositoryTestSuite) Test6_FindByVerificationCode() {
	ctx := context.Background()

	found, err := s.repo.User.FindByVerificationCode(ctx, "123456")
	s.NoError(err)
	s.NotNil(found)
}

func (s *AuthRepositoryTestSuite) Test7_RefreshToken() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	token := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	req := &requests.CreateRefreshToken{
		UserId:    s.userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	res, err := s.repo.RefreshToken.CreateRefreshToken(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(token, res.Token)

	found, err := s.repo.RefreshToken.FindByToken(ctx, token)
	s.NoError(err)
	s.NotNil(found)

	foundByUser, err := s.repo.RefreshToken.FindByUserId(ctx, s.userID)
	s.NoError(err)
	s.NotNil(foundByUser)

	err = s.repo.RefreshToken.DeleteRefreshToken(ctx, token)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) Test8_ResetToken() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	token := "reset-token-123"
	expiresAt := time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05")

	req := &requests.CreateResetTokenRequest{
		UserID:     s.userID,
		ResetToken: token,
		ExpiredAt:  expiresAt,
	}

	res, err := s.repo.ResetToken.CreateResetToken(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(token, res.Token)

	found, err := s.repo.ResetToken.FindByToken(ctx, token)
	s.NoError(err)
	s.NotNil(found)

	err = s.repo.ResetToken.DeleteResetToken(ctx, s.userID)
	s.NoError(err)
}

func TestAuthRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthRepositoryTestSuite))
}
