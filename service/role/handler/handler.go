package handler

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role        pbroles.RoleQueryServiceServer
	RoleCommand pbroles.RoleCommandServiceServer
	UserRole    pbuserrole.UserRoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role:        NewRoleQueryHandleGrpc(deps.Service, deps.Logger),
		RoleCommand: NewRoleCommandHandleGrpc(deps.Service, deps.Logger),
		UserRole:    NewUserRoleHandleGrpc(deps.Service, deps.Logger),
	}
}
