package protomapper

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type categoryProtoMapper struct {
}

func NewCategoryProtoMapper() *categoryProtoMapper {
	return &categoryProtoMapper{}
}

func (c *categoryProtoMapper) ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pbcategories.ApiResponseCategory {
	return &pbcategories.ApiResponseCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pbcategories.ApiResponseCategoryDeleteAt {
	return &pbcategories.ApiResponseCategoryDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryDeleteAt(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pbcategories.ApiResponsesCategory {
	return &pbcategories.ApiResponsesCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDelete(status string, message string) *pbcategories.ApiResponseCategoryDelete {
	return &pbcategories.ApiResponseCategoryDelete{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryAll(status string, message string) *pbcategories.ApiResponseCategoryAll {
	return &pbcategories.ApiResponseCategoryAll{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategoryDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pbcategories.ApiResponsePaginationCategoryDeleteAt {
	return &pbcategories.ApiResponsePaginationCategoryDeleteAt{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategoryDeleteAt(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategory(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pbcategories.ApiResponsePaginationCategory {
	return &pbcategories.ApiResponsePaginationCategory{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategory(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pbcategories.ApiResponseCategoryMonthPrice {
	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryMonthlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pbcategories.ApiResponseCategoryYearPrice {
	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryYearlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pbcategories.ApiResponseCategoryMonthlyTotalPrice {
	return &pbcategories.ApiResponseCategoryMonthlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryMonthlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pbcategories.ApiResponseCategoryYearlyTotalPrice {
	return &pbcategories.ApiResponseCategoryYearlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryYearlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) mapResponseCategory(category *response.CategoryResponse) *pbcategories.CategoryResponse {
	return &pbcategories.CategoryResponse{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategory(categories []*response.CategoryResponse) []*pbcategories.CategoryResponse {
	var mappedCategories []*pbcategories.CategoryResponse

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategory(category))
	}

	return mappedCategories
}

func (c *categoryProtoMapper) mapResponseCategoryDeleteAt(category *response.CategoryResponseDeleteAt) *pbcategories.CategoryResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(*category.DeletedAt)
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategoryDeleteAt(categories []*response.CategoryResponseDeleteAt) []*pbcategories.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategories.CategoryResponseDeleteAt

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategoryDeleteAt(category))
	}

	return mappedCategories
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyPrice(category *response.CategoryMonthPriceResponse) *pbcategories.CategoryMonthPriceResponse {
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   int32(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: int32(category.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryMonthlyPrices(c []*response.CategoryMonthPriceResponse) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyPrice(category *response.CategoryYearPriceResponse) *pbcategories.CategoryYearPriceResponse {
	return &pbcategories.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         int32(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       int32(category.TotalRevenue),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryYearlyPrices(c []*response.CategoryYearPriceResponse) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCashierMonthlyTotalPrice(c *response.CategoriesMonthlyTotalPriceResponse) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyTotalPrices(c []*response.CategoriesMonthlyTotalPriceResponse) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCashierMonthlyTotalPrice(Category))
	}

	return CategoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalSale(c *response.CategoriesYearlyTotalPriceResponse) *pbcategories.CategoriesYearlyTotalPriceResponse {
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalPrices(c []*response.CategoriesYearlyTotalPriceResponse) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCategoryYearlyTotalSale(Category))
	}

	return CategoryRecords
}
