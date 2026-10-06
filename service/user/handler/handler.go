package handler

import (
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-user/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	User        pbusers.UserQueryServiceServer
	UserCommand pbusers.UserCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		User:        NewUserQueryHandleGrpc(deps.Service, deps.Logger),
		UserCommand: NewUserCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
