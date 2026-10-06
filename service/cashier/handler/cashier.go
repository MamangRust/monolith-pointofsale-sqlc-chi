package handler

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

// Map helpers
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

func mapResponseCashier(cashier *models.Cashier) *pbcashiers.CashierResponse {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !cashier.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
	}
	if !cashier.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
	}
	return &pbcashiers.CashierResponse{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
	}
}

func mapResponsesCashier(cashiers []*models.CashierRow) []*pbcashiers.CashierResponse {
	var mappedCashiers []*pbcashiers.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierByMerchant(cashiers []*models.CashierByMerchantRow) []*pbcashiers.CashierResponse {
	var mappedCashiers []*pbcashiers.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponseCashierDeleteAt(cashier *models.Cashier) *pbcashiers.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !cashier.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
	}
	if !cashier.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt.Valid {
		deletedAt = convert.TimeToWrappers(&cashier.DeletedAt.Time)
	}

	return &pbcashiers.CashierResponseDeleteAt{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
		DeletedAt:  deletedAt,
	}
}

func mapResponsesCashierActive(cashiers []*models.CashierActiveRow) []*pbcashiers.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashiers.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt != nil {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierTrashed(cashiers []*models.CashierTrashedRow) []*pbcashiers.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashiers.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&cashier.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt != nil {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponseCashierMonthlySale(cashier *models.CashierMonthlySaleRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySales(c []*models.CashierMonthlySaleRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleById(cashier *models.CashierMonthlySaleByIdRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesById(c []*models.CashierMonthlySaleByIdRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleByMerchant(cashier *models.CashierMonthlySaleByMerchantRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesByMerchant(c []*models.CashierMonthlySaleByMerchantRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySale(cashier *models.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySales(c []*models.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleById(cashier *models.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesById(c []*models.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleByMerchant(cashier *models.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesByMerchant(c []*models.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSales(c []*models.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesById(c []*models.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesByMerchant(c []*models.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSales(c []*models.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesById(c []*models.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesByMerchant(c []*models.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}
