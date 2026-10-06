package pbcompat

// Unified service clients: the flat shared/pb layout exposed one client per
// entity, while the local pb module splits them into Command/Query/Stats
// services. These merged interfaces keep handler code unchanged.
import (
	"google.golang.org/grpc"

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

// CashierServiceClient
type CashierServiceClient interface {
	cashierspb.CashierCommandServiceClient
	cashierspb.CashierQueryServiceClient
	cashierspb.CashierStatsByIdServiceClient
	cashierspb.CashierStatsByMerchantServiceClient
	cashierspb.CashierStatsServiceClient
}

type Cashiercombined struct {
	cashierspb.CashierCommandServiceClient
	cashierspb.CashierQueryServiceClient
	cashierspb.CashierStatsByIdServiceClient
	cashierspb.CashierStatsByMerchantServiceClient
	cashierspb.CashierStatsServiceClient
}

func NewCashierServiceClient(cc grpc.ClientConnInterface) CashierServiceClient {
	return Cashiercombined{
		CashierCommandServiceClient:         cashierspb.NewCashierCommandServiceClient(cc),
		CashierQueryServiceClient:           cashierspb.NewCashierQueryServiceClient(cc),
		CashierStatsByIdServiceClient:       cashierspb.NewCashierStatsByIdServiceClient(cc),
		CashierStatsByMerchantServiceClient: cashierspb.NewCashierStatsByMerchantServiceClient(cc),
		CashierStatsServiceClient:           cashierspb.NewCashierStatsServiceClient(cc),
	}
}

// CategoryServiceClient
type CategoryServiceClient interface {
	categoriespb.CategoryCommandServiceClient
	categoriespb.CategoryQueryServiceClient
	categoriespb.CategoryStatsByIdServiceClient
	categoriespb.CategoryStatsByMerchantServiceClient
	categoriespb.CategoryStatsServiceClient
}

type Categorycombined struct {
	categoriespb.CategoryCommandServiceClient
	categoriespb.CategoryQueryServiceClient
	categoriespb.CategoryStatsByIdServiceClient
	categoriespb.CategoryStatsByMerchantServiceClient
	categoriespb.CategoryStatsServiceClient
}

func NewCategoryServiceClient(cc grpc.ClientConnInterface) CategoryServiceClient {
	return Categorycombined{
		CategoryCommandServiceClient:         categoriespb.NewCategoryCommandServiceClient(cc),
		CategoryQueryServiceClient:           categoriespb.NewCategoryQueryServiceClient(cc),
		CategoryStatsByIdServiceClient:       categoriespb.NewCategoryStatsByIdServiceClient(cc),
		CategoryStatsByMerchantServiceClient: categoriespb.NewCategoryStatsByMerchantServiceClient(cc),
		CategoryStatsServiceClient:           categoriespb.NewCategoryStatsServiceClient(cc),
	}
}

// MerchantServiceClient
type MerchantServiceClient interface {
	merchantspb.MerchantCommandServiceClient
	merchantspb.MerchantQueryServiceClient
}

type Merchantcombined struct {
	merchantspb.MerchantCommandServiceClient
	merchantspb.MerchantQueryServiceClient
}

func NewMerchantServiceClient(cc grpc.ClientConnInterface) MerchantServiceClient {
	return Merchantcombined{
		MerchantCommandServiceClient: merchantspb.NewMerchantCommandServiceClient(cc),
		MerchantQueryServiceClient:   merchantspb.NewMerchantQueryServiceClient(cc),
	}
}

// MerchantDocumentServiceClient
type MerchantDocumentServiceClient interface {
	merchants_documentpb.MerchantDocumentCommandServiceClient
	merchants_documentpb.MerchantDocumentQueryServiceClient
}

type MerchantDocumentcombined struct {
	merchants_documentpb.MerchantDocumentCommandServiceClient
	merchants_documentpb.MerchantDocumentQueryServiceClient
}

func NewMerchantDocumentServiceClient(cc grpc.ClientConnInterface) MerchantDocumentServiceClient {
	return MerchantDocumentcombined{
		MerchantDocumentCommandServiceClient: merchants_documentpb.NewMerchantDocumentCommandServiceClient(cc),
		MerchantDocumentQueryServiceClient:   merchants_documentpb.NewMerchantDocumentQueryServiceClient(cc),
	}
}

// OrderServiceClient
type OrderServiceClient interface {
	orderspb.OrderCommandServiceClient
	orderspb.OrderQueryServiceClient
	orderspb.OrderStatsByIdServiceClient
	orderspb.OrderStatsByMerchantServiceClient
	orderspb.OrderStatsServiceClient
}

type Ordercombined struct {
	orderspb.OrderCommandServiceClient
	orderspb.OrderQueryServiceClient
	orderspb.OrderStatsByIdServiceClient
	orderspb.OrderStatsByMerchantServiceClient
	orderspb.OrderStatsServiceClient
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return Ordercombined{
		OrderCommandServiceClient:         orderspb.NewOrderCommandServiceClient(cc),
		OrderQueryServiceClient:           orderspb.NewOrderQueryServiceClient(cc),
		OrderStatsByIdServiceClient:       orderspb.NewOrderStatsByIdServiceClient(cc),
		OrderStatsByMerchantServiceClient: orderspb.NewOrderStatsByMerchantServiceClient(cc),
		OrderStatsServiceClient:           orderspb.NewOrderStatsServiceClient(cc),
	}
}

// OrderItemServiceClient
type OrderItemServiceClient = order_itemspb.OrderItemQueryServiceClient

var NewOrderItemServiceClient = order_itemspb.NewOrderItemQueryServiceClient

// ProductServiceClient
type ProductServiceClient interface {
	productspb.ProductCommandServiceClient
	productspb.ProductQueryServiceClient
}

type Productcombined struct {
	productspb.ProductCommandServiceClient
	productspb.ProductQueryServiceClient
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return Productcombined{
		ProductCommandServiceClient: productspb.NewProductCommandServiceClient(cc),
		ProductQueryServiceClient:   productspb.NewProductQueryServiceClient(cc),
	}
}

// RoleServiceClient
type RoleServiceClient interface {
	rolespb.RoleCommandServiceClient
	rolespb.RoleQueryServiceClient
}

type Rolecombined struct {
	rolespb.RoleCommandServiceClient
	rolespb.RoleQueryServiceClient
}

func NewRoleServiceClient(cc grpc.ClientConnInterface) RoleServiceClient {
	return Rolecombined{
		RoleCommandServiceClient: rolespb.NewRoleCommandServiceClient(cc),
		RoleQueryServiceClient:   rolespb.NewRoleQueryServiceClient(cc),
	}
}

// TransactionServiceClient
type TransactionServiceClient interface {
	transactionspb.TransactionCommandServiceClient
	transactionspb.TransactionQueryServiceClient
	transactionspb.TransactionStatsByMerchantServiceClient
	transactionspb.TransactionStatsServiceClient
}

type Transactioncombined struct {
	transactionspb.TransactionCommandServiceClient
	transactionspb.TransactionQueryServiceClient
	transactionspb.TransactionStatsByMerchantServiceClient
	transactionspb.TransactionStatsServiceClient
}

func NewTransactionServiceClient(cc grpc.ClientConnInterface) TransactionServiceClient {
	return Transactioncombined{
		TransactionCommandServiceClient:         transactionspb.NewTransactionCommandServiceClient(cc),
		TransactionQueryServiceClient:           transactionspb.NewTransactionQueryServiceClient(cc),
		TransactionStatsByMerchantServiceClient: transactionspb.NewTransactionStatsByMerchantServiceClient(cc),
		TransactionStatsServiceClient:           transactionspb.NewTransactionStatsServiceClient(cc),
	}
}

// UserServiceClient
type UserServiceClient interface {
	userspb.UserCommandServiceClient
	userspb.UserQueryServiceClient
}

type Usercombined struct {
	userspb.UserCommandServiceClient
	userspb.UserQueryServiceClient
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return Usercombined{
		UserCommandServiceClient: userspb.NewUserCommandServiceClient(cc),
		UserQueryServiceClient:   userspb.NewUserQueryServiceClient(cc),
	}
}
