package handler

import (
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
)

type UserQueryHandleGrpc interface {
	pbusers.UserQueryServiceServer
}

type UserCommandHandleGrpc interface {
	pbusers.UserCommandServiceServer
}
