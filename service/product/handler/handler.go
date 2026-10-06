package handler

import (
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Product        pbproducts.ProductQueryServiceServer
	ProductCommand pbproducts.ProductCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product:        NewProductQueryHandleGrpc(deps.Service, deps.Logger),
		ProductCommand: NewProductCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
