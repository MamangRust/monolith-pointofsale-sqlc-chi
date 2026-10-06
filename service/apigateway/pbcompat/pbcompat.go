package pbcompat

// Package pbcompat re-exports the split per-service pb packages under the
// flat names the handlers were written against (former shared/pb layout),
// so handlers can keep referencing every type as pb.X.
import (
	rootpb "github.com/MamangRust/monolith-point-of-sale-pb"
	cashierspb "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	categoriespb "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	merchantspb "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	merchants_documentpb "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	order_itemspb "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	orderspb "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	productspb "github.com/MamangRust/monolith-point-of-sale-pb/products"
	rolespb "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	transactionspb "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	userspb "github.com/MamangRust/monolith-point-of-sale-pb/users"
)

// rootpb
type AuthServiceClient = rootpb.AuthServiceClient
type ForgotPasswordRequest = rootpb.ForgotPasswordRequest
type GetMeRequest = rootpb.GetMeRequest
type LoginRequest = rootpb.LoginRequest

var NewAuthServiceClient = rootpb.NewAuthServiceClient

type RefreshTokenRequest = rootpb.RefreshTokenRequest
type RegisterRequest = rootpb.RegisterRequest
type ResetPasswordRequest = rootpb.ResetPasswordRequest
type VerifyCodeRequest = rootpb.VerifyCodeRequest

// cashierspb
type CreateCashierRequest = cashierspb.CreateCashierRequest
type FindAllCashierRequest = cashierspb.FindAllCashierRequest
type FindByIdCashierRequest = cashierspb.FindByIdCashierRequest
type FindYearCashier = cashierspb.FindYearCashier
type FindYearCashierById = cashierspb.FindYearCashierById
type FindYearCashierByMerchant = cashierspb.FindYearCashierByMerchant
type FindYearMonthTotalSales = cashierspb.FindYearMonthTotalSales
type FindYearMonthTotalSalesById = cashierspb.FindYearMonthTotalSalesById
type FindYearMonthTotalSalesByMerchant = cashierspb.FindYearMonthTotalSalesByMerchant
type FindYearTotalSales = cashierspb.FindYearTotalSales
type FindYearTotalSalesById = cashierspb.FindYearTotalSalesById
type FindYearTotalSalesByMerchant = cashierspb.FindYearTotalSalesByMerchant
type UpdateCashierRequest = cashierspb.UpdateCashierRequest

// categoriespb
type CreateCategoryRequest = categoriespb.CreateCategoryRequest
type FindAllCategoryRequest = categoriespb.FindAllCategoryRequest
type FindByIdCategoryRequest = categoriespb.FindByIdCategoryRequest
type FindYearCategory = categoriespb.FindYearCategory
type FindYearCategoryById = categoriespb.FindYearCategoryById
type FindYearCategoryByMerchant = categoriespb.FindYearCategoryByMerchant
type FindYearMonthTotalPriceById = categoriespb.FindYearMonthTotalPriceById
type FindYearMonthTotalPriceByMerchant = categoriespb.FindYearMonthTotalPriceByMerchant
type FindYearMonthTotalPrices = categoriespb.FindYearMonthTotalPrices
type FindYearTotalPriceById = categoriespb.FindYearTotalPriceById
type FindYearTotalPriceByMerchant = categoriespb.FindYearTotalPriceByMerchant
type FindYearTotalPrices = categoriespb.FindYearTotalPrices
type UpdateCategoryRequest = categoriespb.UpdateCategoryRequest

// merchantspb
type CreateMerchantRequest = merchantspb.CreateMerchantRequest
type FindAllMerchantRequest = merchantspb.FindAllMerchantRequest
type FindByIdMerchantRequest = merchantspb.FindByIdMerchantRequest
type UpdateMerchantRequest = merchantspb.UpdateMerchantRequest
type UpdateMerchantStatusRequest = merchantspb.UpdateMerchantStatusRequest

// merchants_documentpb
type CreateMerchantDocumentRequest = merchants_documentpb.CreateMerchantDocumentRequest
type DeleteMerchantDocumentPermanentRequest = merchants_documentpb.DeleteMerchantDocumentPermanentRequest
type FindAllMerchantDocumentsRequest = merchants_documentpb.FindAllMerchantDocumentsRequest
type FindMerchantDocumentByIdRequest = merchants_documentpb.FindMerchantDocumentByIdRequest
type RestoreMerchantDocumentRequest = merchants_documentpb.RestoreMerchantDocumentRequest
type TrashedMerchantDocumentRequest = merchants_documentpb.TrashedMerchantDocumentRequest
type UpdateMerchantDocumentRequest = merchants_documentpb.UpdateMerchantDocumentRequest
type UpdateMerchantDocumentStatusRequest = merchants_documentpb.UpdateMerchantDocumentStatusRequest

// order_itemspb
type FindAllOrderItemRequest = order_itemspb.FindAllOrderItemRequest
type FindByIdOrderItemRequest = order_itemspb.FindByIdOrderItemRequest

// orderspb
type CreateOrderItemRequest = orderspb.CreateOrderItemRequest
type CreateOrderRequest = orderspb.CreateOrderRequest
type FindAllOrderRequest = orderspb.FindAllOrderRequest
type FindByIdOrderRequest = orderspb.FindByIdOrderRequest
type FindYearMonthTotalRevenue = orderspb.FindYearMonthTotalRevenue
type FindYearMonthTotalRevenueByMerchant = orderspb.FindYearMonthTotalRevenueByMerchant
type FindYearOrder = orderspb.FindYearOrder
type FindYearOrderByMerchant = orderspb.FindYearOrderByMerchant
type FindYearTotalRevenue = orderspb.FindYearTotalRevenue
type FindYearTotalRevenueByMerchant = orderspb.FindYearTotalRevenueByMerchant
type UpdateOrderItemRequest = orderspb.UpdateOrderItemRequest
type UpdateOrderRequest = orderspb.UpdateOrderRequest

// productspb
type CreateProductRequest = productspb.CreateProductRequest
type FindAllProductCategoryRequest = productspb.FindAllProductCategoryRequest
type FindAllProductMerchantRequest = productspb.FindAllProductMerchantRequest
type FindAllProductRequest = productspb.FindAllProductRequest
type FindByIdProductRequest = productspb.FindByIdProductRequest
type UpdateProductRequest = productspb.UpdateProductRequest

// rolespb
type CreateRoleRequest = rolespb.CreateRoleRequest
type FindAllRoleRequest = rolespb.FindAllRoleRequest
type FindByIdRoleRequest = rolespb.FindByIdRoleRequest
type FindByIdUserRoleRequest = rolespb.FindByIdUserRoleRequest
type UpdateRoleRequest = rolespb.UpdateRoleRequest

// transactionspb
type CreateTransactionRequest = transactionspb.CreateTransactionRequest
type FindAllTransactionMerchantRequest = transactionspb.FindAllTransactionMerchantRequest
type FindAllTransactionRequest = transactionspb.FindAllTransactionRequest
type FindByIdTransactionRequest = transactionspb.FindByIdTransactionRequest
type FindMonthlyTransactionStatus = transactionspb.FindMonthlyTransactionStatus
type FindMonthlyTransactionStatusByMerchant = transactionspb.FindMonthlyTransactionStatusByMerchant
type FindYearlyTransactionStatus = transactionspb.FindYearlyTransactionStatus
type FindYearlyTransactionStatusByMerchant = transactionspb.FindYearlyTransactionStatusByMerchant
type MonthTransactionMethod = transactionspb.MonthTransactionMethod
type MonthTransactionMethodByMerchant = transactionspb.MonthTransactionMethodByMerchant
type UpdateTransactionRequest = transactionspb.UpdateTransactionRequest
type YearTransactionMethod = transactionspb.YearTransactionMethod
type YearTransactionMethodByMerchant = transactionspb.YearTransactionMethodByMerchant

// userspb
type CreateUserRequest = userspb.CreateUserRequest
type FindAllUserRequest = userspb.FindAllUserRequest
type FindByIdUserRequest = userspb.FindByIdUserRequest
type UpdateUserRequest = userspb.UpdateUserRequest
