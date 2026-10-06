package response_api

import (
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type orderResponseMapper struct {
}

func NewOrderResponseMapper() *orderResponseMapper {
	return &orderResponseMapper{}
}

func (o *orderResponseMapper) ToResponseOrder(order *pborders.OrderResponse) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		CashierID:  int(order.CashierId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrder(orders []*pborders.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderResponseMapper) ToResponseOrderDeleteAt(order *pborders.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
	var deletedAt string

	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &response.OrderResponseDeleteAt{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		CashierID:  int(order.CashierId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeleteAt:   &deletedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrderDeleteAt(orders []*pborders.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (s *orderResponseMapper) ToOrderMonthlyPrice(category *pborders.OrderMonthlyResponse) *response.OrderMonthlyResponse {
	return &response.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int(category.OrderCount),
		TotalRevenue:   int(category.TotalRevenue),
		TotalItemsSold: int(category.TotalItemsSold),
	}
}

func (s *orderResponseMapper) ToOrderMonthlyPrices(c []*pborders.OrderMonthlyResponse) []*response.OrderMonthlyResponse {
	var categoryRecords []*response.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToOrderYearlyPrice(category *pborders.OrderYearlyResponse) *response.OrderYearlyResponse {
	return &response.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int(category.OrderCount),
		TotalRevenue:       int(category.TotalRevenue),
		TotalItemsSold:     int(category.TotalItemsSold),
		ActiveCashiers:     int(category.ActiveCashiers),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *orderResponseMapper) ToOrderYearlyPrices(c []*pborders.OrderYearlyResponse) []*response.OrderYearlyResponse {
	var categoryRecords []*response.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToResponseOrderMonthlyTotalRevenue(c *pborders.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse {
	return &response.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		TotalRevenue:   int(c.TotalRevenue),
		TotalItemsSold: int(c.TotalItemsSold),
	}
}

func (s *orderResponseMapper) ToResponseOrderMonthlyTotalRevenues(c []*pborders.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*response.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToResponseOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderResponseMapper) ToResponseOrderYearlyTotalRevenue(c *pborders.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse {
	return &response.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderResponseMapper) ToResponseOrderYearlyTotalRevenues(c []*pborders.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse {
	var orderRecords []*response.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToResponseOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}

func (o *orderResponseMapper) ToApiResponseOrder(pbResponse *pborders.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pborders.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponsesOrder(pbResponse *pborders.ApiResponsesOrder) *response.ApiResponsesOrder {
	return &response.ApiResponsesOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponsesOrder(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDelete(pbResponse *pborders.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponseOrderAll(pbResponse *pborders.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrderDeleteAt(pbResponse *pborders.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt {
	return &response.ApiResponsePaginationOrderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrder(pbResponse *pborders.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder {
	return &response.ApiResponsePaginationOrder{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrder(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderResponseMapper) ToApiResponseMonthlyOrder(pbResponse *pborders.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly {
	return &response.ApiResponseOrderMonthly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToOrderMonthlyPrices(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseYearlyOrder(pbResponse *pborders.ApiResponseOrderYearly) *response.ApiResponseOrderYearly {
	return &response.ApiResponseOrderYearly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToOrderYearlyPrices(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseMonthlyTotalRevenue(pbResponse *pborders.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue {
	return &response.ApiResponseOrderMonthlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderMonthlyTotalRevenues(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseYearlyTotalRevenue(pbResponse *pborders.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue {
	return &response.ApiResponseOrderYearlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderYearlyTotalRevenues(pbResponse.Data),
	}
}
