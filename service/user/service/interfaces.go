package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, error)
}

type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error)
	CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*models.User, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error)
	TrashedUser(ctx context.Context, user_id int) (*models.User, error)
	RestoreUser(ctx context.Context, user_id int) (*models.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}
