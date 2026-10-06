package handler

import (
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
)

type ProductQueryHandleGrpc interface {
	pbproducts.ProductQueryServiceServer
}

type ProductCommandHandleGrpc interface {
	pbproducts.ProductCommandServiceServer
}
