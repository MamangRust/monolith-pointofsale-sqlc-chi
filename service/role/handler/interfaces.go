package handler

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
)

type RoleQueryHandleGrpc interface {
	pbroles.RoleQueryServiceServer
}

type RoleCommandHandleGrpc interface {
	pbroles.RoleCommandServiceServer
}
