package protomapper

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cashierProtoMapper struct {
}

func NewCashierProtoMapper() *cashierProtoMapper {
	return &cashierProtoMapper{}
}

func (c *cashierProtoMapper) ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pbcashiers.ApiResponseCashier {
	return &pbcashiers.ApiResponseCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pbcashiers.ApiResponsesCashier {
	return &pbcashiers.ApiResponsesCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pbcashiers.ApiResponseCashierDeleteAt {
	return &pbcashiers.ApiResponseCashierDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashierDeleteAt(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDelete(status string, message string) *pbcashiers.ApiResponseCashierDelete {
	return &pbcashiers.ApiResponseCashierDelete{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponseCashierAll(status string, message string) *pbcashiers.ApiResponseCashierAll {
	return &pbcashiers.ApiResponseCashierAll{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pbcashiers.ApiResponsePaginationCashierDeleteAt {
	return &pbcashiers.ApiResponsePaginationCashierDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashierDeleteAt(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pbcashiers.ApiResponsePaginationCashier {
	return &pbcashiers.ApiResponsePaginationCashier{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashier(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbcashiers.ApiResponseCashierMonthSales {
	return &pbcashiers.ApiResponseCashierMonthSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierMonthlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbcashiers.ApiResponseCashierYearSales {
	return &pbcashiers.ApiResponseCashierYearSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierYearlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbcashiers.ApiResponseCashierMonthlyTotalSales {
	return &pbcashiers.ApiResponseCashierMonthlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierMonthlyTotalSales(row),
	}
}

func (u *cashierProtoMapper) ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbcashiers.ApiResponseCashierYearlyTotalSales {
	return &pbcashiers.ApiResponseCashierYearlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierYearlyTotalSales(row),
	}
}

func (c *cashierProtoMapper) mapResponseCashier(cashier *response.CashierResponse) *pbcashiers.CashierResponse {
	return &pbcashiers.CashierResponse{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashier(cashiers []*response.CashierResponse) []*pbcashiers.CashierResponse {
	var mappedCashiers []*pbcashiers.CashierResponse

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashier(cashier))
	}

	return mappedCashiers
}

func (c *cashierProtoMapper) mapResponseCashierDeleteAt(cashier *response.CashierResponseDeleteAt) *pbcashiers.CashierResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt != nil {
		deletedAt = wrapperspb.String(*cashier.DeletedAt)
	}

	return &pbcashiers.CashierResponseDeleteAt{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashierDeleteAt(cashiers []*response.CashierResponseDeleteAt) []*pbcashiers.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashiers.CashierResponseDeleteAt

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashierDeleteAt(cashier))
	}

	return mappedCashiers
}

func (s *cashierProtoMapper) mapResponseCashierMonthlySale(cashier *response.CashierResponseMonthSales) *pbcashiers.CashierResponseMonthSales {
	return &pbcashiers.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierMonthlySales(c []*response.CashierResponseMonthSales) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlySale(cashier *response.CashierResponseYearSales) *pbcashiers.CashierResponseYearSales {
	return &pbcashiers.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierYearlySales(c []*response.CashierResponseYearSales) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSale(c *response.CashierResponseMonthTotalSales) *pbcashiers.CashierResponseMonthTotalSales {
	return &pbcashiers.CashierResponseMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSales(c []*response.CashierResponseMonthTotalSales) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSale(c *response.CashierResponseYearTotalSales) *pbcashiers.CashierResponseYearTotalSales {
	return &pbcashiers.CashierResponseYearTotalSales{
		Year:       c.Year,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSales(c []*response.CashierResponseYearTotalSales) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}
