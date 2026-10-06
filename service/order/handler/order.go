package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseOrder(order *models.Order) *pborders.OrderResponse {
	if order == nil {
		return nil
	}
	return &pborders.OrderResponse{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  convert.FormatTimePtr(&order.CreatedAt),
		UpdatedAt:  convert.FormatTimePtr(&order.UpdatedAt),
	}
}

func mapResponsesOrder(orders []*models.OrderRow) []*pborders.OrderResponse {
	var mappedOrders []*pborders.OrderResponse
	for _, order := range orders {
		if order == nil {
			continue
		}
		mappedOrders = append(mappedOrders, &pborders.OrderResponse{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  convert.FormatTimePtr(&order.CreatedAt),
			UpdatedAt:  convert.FormatTimePtr(&order.UpdatedAt),
		})
	}
	return mappedOrders
}

func mapResponseOrderDeleteAt(order *models.Order) *pborders.OrderResponseDeleteAt {
	if order == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if order.DeletedAt.Valid {
		deletedAt = convert.TimeToWrappers(&order.DeletedAt.Time)
	}

	return &pborders.OrderResponseDeleteAt{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  convert.FormatTimePtr(&order.CreatedAt),
		UpdatedAt:  convert.FormatTimePtr(&order.UpdatedAt),
		DeletedAt:  deletedAt,
	}
}

func mapResponsesOrderActive(orders []*models.OrderActiveRow) []*pborders.OrderResponseDeleteAt {
	var mappedOrders []*pborders.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt != nil {
			deletedAt = wrapperspb.String(order.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pborders.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  convert.FormatTimePtr(&order.CreatedAt),
			UpdatedAt:  convert.FormatTimePtr(&order.UpdatedAt),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponsesOrderTrashed(orders []*models.OrderTrashedRow) []*pborders.OrderResponseDeleteAt {
	var mappedOrders []*pborders.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt != nil {
			deletedAt = wrapperspb.String(order.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pborders.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  convert.FormatTimePtr(&order.CreatedAt),
			UpdatedAt:  convert.FormatTimePtr(&order.UpdatedAt),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponseOrderMonthlyTotalRevenue(row *models.MonthlyTotalRevenueRow) *pborders.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborders.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenues(c []*models.MonthlyTotalRevenueRow) []*pborders.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderMonthlyTotalRevenueByMerchant(row *models.MonthlyTotalRevenueByMerchantRow) *pborders.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborders.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenuesByMerchant(c []*models.MonthlyTotalRevenueByMerchantRow) []*pborders.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenue(row *models.YearlyTotalRevenueRow) *pborders.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborders.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenues(c []*models.YearlyTotalRevenueRow) []*pborders.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenueByMerchant(row *models.YearlyTotalRevenueByMerchantRow) *pborders.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborders.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenuesByMerchant(c []*models.YearlyTotalRevenueByMerchantRow) []*pborders.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborders.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponsesOrderMonthlyPrices(c []*models.MonthlyOrderRow) []*pborders.OrderMonthlyResponse {
	var categoryRecords []*pborders.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborders.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderMonthlyPricesByMerchant(c []*models.MonthlyOrderByMerchantRow) []*pborders.OrderMonthlyResponse {
	var categoryRecords []*pborders.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborders.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPrices(c []*models.YearlyOrderRow) []*pborders.OrderYearlyResponse {
	var categoryRecords []*pborders.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborders.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPricesByMerchant(c []*models.YearlyOrderByMerchantRow) []*pborders.OrderYearlyResponse {
	var categoryRecords []*pborders.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborders.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}
