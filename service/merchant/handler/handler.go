package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Merchant                pbmerchants.MerchantQueryServiceServer
	MerchantCommand         pbmerchants.MerchantCommandServiceServer
	MerchantDocument        pbmerchants_document.MerchantDocumentQueryServiceServer
	MerchantDocumentCommand pbmerchants_document.MerchantDocumentCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Merchant:                NewMerchantQueryHandleGrpc(deps.Service, deps.Logger),
		MerchantCommand:         NewMerchantCommandHandleGrpc(deps.Service, deps.Logger),
		MerchantDocument:        NewMerchantDocumentQueryHandleGrpc(deps.Service, deps.Logger),
		MerchantDocumentCommand: NewMerchantDocumentCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
