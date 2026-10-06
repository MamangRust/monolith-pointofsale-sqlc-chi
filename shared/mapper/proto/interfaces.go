package protomapper

import (
	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type AuthProtoMapper interface {
	ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode
	ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword
	ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbusers.ApiResponseUserDeleteAt
	ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pbusers.ApiResponsesUser
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbusers.ApiResponseUser
	ToProtoResponseUserDelete(status string, message string) *pbusers.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pbusers.ApiResponseUserAll
	ToProtoResponsePaginationUserDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pbusers.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponse) *pbusers.ApiResponsePaginationUser
}

type RoleProtoMapper interface {
	ToProtoResponseRoleAll(status string, message string) *pbroles.ApiResponseRoleAll
	ToProtoResponseRoleDelete(status string, message string) *pbroles.ApiResponseRoleDelete
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbroles.ApiResponseRole
	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbroles.ApiResponsesRole
	ToProtoResponsePaginationRole(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pbroles.ApiResponsePaginationRole
	ToProtoResponsePaginationRoleDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbroles.ApiResponsePaginationRoleDeleteAt
}

type CategoryProtoMapper interface {
	ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pbcategories.ApiResponseCategoryMonthlyTotalPrice
	ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pbcategories.ApiResponseCategoryYearlyTotalPrice
	ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pbcategories.ApiResponseCategoryMonthPrice
	ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pbcategories.ApiResponseCategoryYearPrice

	ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pbcategories.ApiResponsesCategory
	ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pbcategories.ApiResponseCategoryDeleteAt

	ToProtoResponseCategoryAll(status string, message string) *pbcategories.ApiResponseCategoryAll
	ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pbcategories.ApiResponseCategory
	ToProtoResponseCategoryDelete(status string, message string) *pbcategories.ApiResponseCategoryDelete
	ToProtoResponsePaginationCategoryDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pbcategories.ApiResponsePaginationCategoryDeleteAt
	ToProtoResponsePaginationCategory(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pbcategories.ApiResponsePaginationCategory
}

type CashierProtoMapper interface {
	ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbcashiers.ApiResponseCashierMonthlyTotalSales
	ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbcashiers.ApiResponseCashierYearlyTotalSales

	ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbcashiers.ApiResponseCashierMonthSales
	ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbcashiers.ApiResponseCashierYearSales

	ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pbcashiers.ApiResponseCashier
	ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pbcashiers.ApiResponseCashierDeleteAt
	ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pbcashiers.ApiResponsesCashier
	ToProtoResponseCashierDelete(status string, message string) *pbcashiers.ApiResponseCashierDelete
	ToProtoResponseCashierAll(status string, message string) *pbcashiers.ApiResponseCashierAll
	ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pbcashiers.ApiResponsePaginationCashierDeleteAt
	ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pbcashiers.ApiResponsePaginationCashier
}

type MerchantProtoMapper interface {
	ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pbmerchants.ApiResponseMerchant
	ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pbmerchants.ApiResponseMerchantDeleteAt

	ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pbmerchants.ApiResponsesMerchant
	ToProtoResponseMerchantDelete(status string, message string) *pbmerchants.ApiResponseMerchantDelete
	ToProtoResponseMerchantAll(status string, message string) *pbmerchants.ApiResponseMerchantAll
	ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pbmerchants.ApiResponsePaginationMerchantDeleteAt
	ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pbmerchants.ApiResponsePaginationMerchant
}

type MerchantDocumentProtoMapper interface {
	ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pbmerchants_document.ApiResponseMerchantDocument
	ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchants_document.ApiResponsesMerchantDocument

	ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchants_document.ApiResponsePaginationMerchantDocument
	ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pbmerchants_document.ApiResponsePaginationMerchantDocumentAt

	ToProtoResponseMerchantDocumentDelete(status string, message string) *pbmerchants_document.ApiResponseMerchantDocumentDelete

	ToProtoResponseMerchantDocumentAll(status string, message string) *pbmerchants_document.ApiResponseMerchantDocumentAll
}

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pborder_items.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pborder_items.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pborder_items.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pborder_items.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pborder_items.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pborder_items.ApiResponsePaginationOrderItem
}

type OrderProtoMapper interface {
	ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pborders.ApiResponseOrderMonthlyTotalRevenue
	ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pborders.ApiResponseOrderYearlyTotalRevenue

	ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pborders.ApiResponseOrderMonthly
	ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pborders.ApiResponseOrderYearly

	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborders.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborders.ApiResponseOrderDeleteAt
	ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pborders.ApiResponsesOrder
	ToProtoResponseOrderDelete(status string, message string) *pborders.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pborders.ApiResponseOrderAll
	ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pborders.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pborders.ApiResponsePaginationOrder
}

type ProductProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproducts.ApiResponseProduct
	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproducts.ApiResponseProductDeleteAt

	ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pbproducts.ApiResponsesProduct
	ToProtoResponseProductDelete(status string, message string) *pbproducts.ApiResponseProductDelete
	ToProtoResponseProductAll(status string, message string) *pbproducts.ApiResponseProductAll
	ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pbproducts.ApiResponsePaginationProductDeleteAt
	ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pbproducts.ApiResponsePaginationProduct
}

type TransactionProtoMapper interface {
	ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbtransactions.ApiResponseTransactionMonthAmountSuccess
	ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbtransactions.ApiResponseTransactionYearAmountSuccess
	ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbtransactions.ApiResponseTransactionMonthAmountFailed
	ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbtransactions.ApiResponseTransactionYearAmountFailed
	ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbtransactions.ApiResponseTransactionMonthPaymentMethod
	ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbtransactions.ApiResponseTransactionYearPaymentmethod

	ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pbtransactions.ApiResponseTransaction
	ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pbtransactions.ApiResponseTransactionDeleteAt
	ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pbtransactions.ApiResponsesTransaction
	ToProtoResponseTransactionDelete(status string, message string) *pbtransactions.ApiResponseTransactionDelete
	ToProtoResponseTransactionAll(status string, message string) *pbtransactions.ApiResponseTransactionAll
	ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pbtransactions.ApiResponsePaginationTransactionDeleteAt
	ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pbtransactions.ApiResponsePaginationTransaction
}
