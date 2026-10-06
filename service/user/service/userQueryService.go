package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userQueryDeps struct {
	Cache         mencache.UserQueryCache
	UserQuery     repository.UserQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type userQueryService struct {
	mencache      mencache.UserQueryCache
	userQuery     repository.UserQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewUserQueryService(params *userQueryDeps) UserQueryService {
	return &userQueryService{
		mencache:      params.Cache,
		userQuery:     params.UserQuery,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully fetched users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.userQuery.FindAllUsers(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*models.UserRow](
			s.logger,
			user_errors.ErrFailedFindAll.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedUsersCache(ctx, req, res, &totalCount)
	logSuccess("Successfully fetched users", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return res, &totalCount, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*models.User, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully fetched user from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userQuery.FindById(ctx, id)
	if err != nil {
		status = "error"
		// Propagate the typed repository error unchanged: ErrUserNotFound
		// (404) for no-rows, ErrInternalServerError (500) for other DB
		// failures — do not mask them into a generic 404.
		return sharederrorhandler.HandleError[*models.User](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", id),
			zap.Error(err),
		)
	}

	s.mencache.SetCachedUserCache(ctx, user)
	logSuccess("Successfully fetched user", zap.Int("user.id", id))
	return user, nil
}

// FindByEmailWithPassword resolves a user (including the password hash) by
// email for credential verification. It is never cached.
func (s *userQueryService) FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	const method = "FindByEmailWithPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("user.email", email),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByEmailWithPassword(ctx, email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.Error(err))
	}
	if user == nil {
		status = "error"
		return nil, user_errors.ErrUserNotFound
	}

	logSuccess("Successfully fetched user by email", zap.String("user.email", email))
	return user, nil
}

// FindByEmailAndVerify resolves a verified user by email for login.
func (s *userQueryService) FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error) {
	const method = "FindByEmailAndVerify"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("user.email", email),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByEmailAndVerify(ctx, email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.Error(err))
	}
	if user == nil {
		status = "error"
		return nil, user_errors.ErrUserNotFound
	}

	logSuccess("Successfully fetched verified user by email", zap.String("user.email", email))
	return user, nil
}

// FindByVerificationCode resolves a user by their verification code.
func (s *userQueryService) FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error) {
	const method = "FindByVerificationCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByVerificationCode(ctx, verification_code)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.Error(err))
	}
	if user == nil {
		status = "error"
		return nil, user_errors.ErrUserNotFound
	}

	logSuccess("Successfully fetched user by verification code")
	return user, nil
}

func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.userQuery.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*models.UserActiveRow](
			s.logger,
			user_errors.ErrFailedFindActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedUserActiveCache(ctx, req, res, &totalCount)
	logSuccess("Successfully fetched active users", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return res, &totalCount, nil
}

func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.userQuery.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*models.UserTrashedRow](
			s.logger,
			user_errors.ErrFailedFindTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedUserTrashedCache(ctx, req, res, &totalCount)
	logSuccess("Successfully fetched trashed users", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return res, &totalCount, nil
}

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
