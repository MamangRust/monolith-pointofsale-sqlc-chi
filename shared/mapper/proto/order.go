package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderProtoMapper struct{}

func NewOrderProtoMapper() *orderProtoMapper {
	return &orderProtoMapper{}
}

func (o *orderProtoMapper) ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborders.ApiResponseOrder {
	return &pborders.ApiResponseOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pborders.ApiResponsesOrder {
	return &pborders.ApiResponsesOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborders.ApiResponseOrderDeleteAt {
	return &pborders.ApiResponseOrderDeleteAt{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderDeleteAt(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDelete(status string, message string) *pborders.ApiResponseOrderDelete {
	return &pborders.ApiResponseOrderDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderAll(status string, message string) *pborders.ApiResponseOrderAll {
	return &pborders.ApiResponseOrderAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pborders.ApiResponsePaginationOrderDeleteAt {
	return &pborders.ApiResponsePaginationOrderDeleteAt{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderDeleteAt(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pborders.ApiResponsePaginationOrder {
	return &pborders.ApiResponsePaginationOrder{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrder(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pborders.ApiResponseOrderMonthly {
	return &pborders.ApiResponseOrderMonthly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderMonthlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pborders.ApiResponseOrderYearly {
	return &pborders.ApiResponseOrderYearly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderYearlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pborders.ApiResponseOrderMonthlyTotalRevenue {
	return &pborders.ApiResponseOrderMonthlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderMonthlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pborders.ApiResponseOrderYearlyTotalRevenue {
	return &pborders.ApiResponseOrderYearlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderYearlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) mapResponseOrder(order *response.OrderResponse) *pborders.OrderResponse {
	return &pborders.OrderResponse{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrder(orders []*response.OrderResponse) []*pborders.OrderResponse {
	var mappedOrders []*pborders.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderProtoMapper) mapResponseOrderDeleteAt(order *response.OrderResponseDeleteAt) *pborders.OrderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if order.DeleteAt != nil {
		deletedAt = wrapperspb.String(*order.DeleteAt)
	}

	return &pborders.OrderResponseDeleteAt{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrderDeleteAt(orders []*response.OrderResponseDeleteAt) []*pborders.OrderResponseDeleteAt {
	var mappedOrders []*pborders.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (s *orderProtoMapper) mapResponseOrderMonthlyPrice(category *response.OrderMonthlyResponse) *pborders.OrderMonthlyResponse {
	return &pborders.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int32(category.OrderCount),
		TotalRevenue:   int32(category.TotalRevenue),
		TotalItemsSold: int32(category.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderMonthlyPrices(c []*response.OrderMonthlyResponse) []*pborders.OrderMonthlyResponse {
	var categoryRecords []*pborders.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyPrice(category *response.OrderYearlyResponse) *pborders.OrderYearlyResponse {
	return &pborders.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int32(category.OrderCount),
		TotalRevenue:       int32(category.TotalRevenue),
		TotalItemsSold:     int32(category.TotalItemsSold),
		ActiveCashiers:     int32(category.ActiveCashiers),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderYearlyPrices(c []*response.OrderYearlyResponse) []*pborders.OrderYearlyResponse {
	var categoryRecords []*pborders.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenue(c *response.OrderMonthlyTotalRevenueResponse) *pborders.OrderMonthlyTotalRevenueResponse {
	return &pborders.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		TotalRevenue:   int32(c.TotalRevenue),
		TotalItemsSold: int32(c.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenues(c []*response.OrderMonthlyTotalRevenueResponse) []*pborders.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenue(c *response.OrderYearlyTotalRevenueResponse) *pborders.OrderYearlyTotalRevenueResponse {
	return &pborders.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenues(c []*response.OrderYearlyTotalRevenueResponse) []*pborders.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}
