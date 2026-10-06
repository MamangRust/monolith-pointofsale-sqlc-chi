package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserRow, total *int)

	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserActiveRow, total *int)

	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserTrashedRow, total *int)

	GetCachedUserCache(ctx context.Context, id int) (*models.User, bool)
	SetCachedUserCache(ctx context.Context, data *models.User)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
	DeleteUserAllCache(ctx context.Context)
}
