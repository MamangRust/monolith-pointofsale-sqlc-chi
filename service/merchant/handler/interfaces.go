package handler

import (
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
)

type MerchantDocumentQueryHandleGrpc interface {
	pbmerchants_document.MerchantDocumentQueryServiceServer
}

type MerchantDocumentCommandHandleGrpc interface {
	pbmerchants_document.MerchantDocumentCommandServiceServer
}
