package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	auth_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/auth"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	pb "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-pkg/upload_image"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
)

type ServiceConnections struct {
	Auth        *grpc.ClientConn
	Role        *grpc.ClientConn
	User        *grpc.ClientConn
	Cashier     *grpc.ClientConn
	Category    *grpc.ClientConn
	Merchant    *grpc.ClientConn
	OrderItem   *grpc.ClientConn
	Order       *grpc.ClientConn
	Product     *grpc.ClientConn
	Transaction *grpc.ClientConn
}

type Deps struct {
	Token              auth.TokenManager
	Router             chi.Router
	Logger             logger.LoggerInterface
	Mapping            *response_api.ResponseApiMapper
	ImageUpload        upload_image.ImageUploads
	ServiceConnections *ServiceConnections
	ApiHandler         apierror.ApiHandler
	AuthCache          auth_cache.AuthMencache
	GatewayCache       *gateway_cache.GatewayCache
}

func NewHandler(deps *Deps) {
	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pb.NewRoleServiceClient(deps.ServiceConnections.Role)
	clientUser := pb.NewUserServiceClient(deps.ServiceConnections.User)
	clientCategory := pb.NewCategoryServiceClient(deps.ServiceConnections.Category)
	clientCashier := pb.NewCashierServiceClient(deps.ServiceConnections.Cashier)
	clientMerchant := pb.NewMerchantServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantDocument := pb.NewMerchantDocumentServiceClient(deps.ServiceConnections.Merchant)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.ServiceConnections.OrderItem)
	clientOrder := pb.NewOrderServiceClient(deps.ServiceConnections.Order)
	clientProduct := pb.NewProductServiceClient(deps.ServiceConnections.Product)
	clientTransaction := pb.NewTransactionServiceClient(deps.ServiceConnections.Transaction)

	NewHandlerAuth(deps.Router, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper, deps.ApiHandler, deps.AuthCache)
	NewHandlerRole(deps.Router, clientRole, deps.Logger, deps.Mapping.RoleResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerUser(deps.Router, clientUser, deps.Logger, deps.Mapping.UserResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerCategory(deps.Router, clientCategory, deps.Logger, deps.Mapping.CategoryResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerCashier(deps.Router, clientCashier, deps.Logger, deps.Mapping.CashierResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerMerchant(deps.Router, clientMerchant, deps.Logger, deps.Mapping.MerchantResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerMerchantDocument(deps.Router, clientMerchantDocument, deps.Logger, deps.Mapping.MerchantDocumentProMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerOrderItem(deps.Router, clientOrderItem, deps.Logger, deps.Mapping.OrderItemResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerOrder(deps.Router, clientOrder, deps.Logger, deps.Mapping.OrderResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerProduct(deps.Router, clientProduct, deps.Logger, deps.Mapping.ProductResponseMapper, deps.ImageUpload, deps.ApiHandler, deps.GatewayCache)
	NewHandlerTransaction(deps.Router, clientTransaction, deps.Logger, deps.Mapping.TransactionResponseMapper, deps.ApiHandler, deps.GatewayCache)
}

func parseQueryInt(r *http.Request, key string, defaultValue int) int {
	val, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func parseQueryIntWithValidation(r *http.Request, key string, min, max int) (int, error) {
	valStr := r.URL.Query().Get(key)
	val, err := strconv.Atoi(valStr)
	if err != nil || val < min || val > max {
		return 0, fmt.Errorf("invalid %s: %s", key, valStr)
	}
	return val, nil
}
