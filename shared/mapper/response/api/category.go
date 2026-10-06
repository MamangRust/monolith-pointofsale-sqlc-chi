package response_api

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type categoryResponseMapper struct{}

func NewCategoryResponseMapper() *categoryResponseMapper {
	return &categoryResponseMapper{}
}

func (c *categoryResponseMapper) ToResponseCategory(category *pbcategories.CategoryResponse) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryResponseMapper) ToResponsesCategory(categories []*pbcategories.CategoryResponse) []*response.CategoryResponse {
	var mappedCategories []*response.CategoryResponse

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategory(category))
	}

	return mappedCategories
}

func (c *categoryResponseMapper) ToResponseCategoryDelete(category *pbcategories.CategoryResponseDeleteAt) *response.CategoryResponseDeleteAt {
	var deletedAt string
	if category.DeletedAt != nil {
		deletedAt = category.DeletedAt.Value
	}

	return &response.CategoryResponseDeleteAt{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (s *categoryResponseMapper) ToResponseCategoryMonthlyPrice(category *pbcategories.CategoryMonthPriceResponse) *response.CategoryMonthPriceResponse {
	return &response.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryID:   int(category.CategoryId),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToResponseCategoryMonthlyPrices(c []*pbcategories.CategoryMonthPriceResponse) []*response.CategoryMonthPriceResponse {
	var categoryRecords []*response.CategoryMonthPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToResponseCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryResponseMapper) ToResponseCategoryYearlyPrice(category *pbcategories.CategoryYearPriceResponse) *response.CategoryYearPriceResponse {
	return &response.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryID:         int(category.CategoryId),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *categoryResponseMapper) ToResponseCategoryYearlyPrices(c []*pbcategories.CategoryYearPriceResponse) []*response.CategoryYearPriceResponse {
	var categoryRecords []*response.CategoryYearPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToResponseCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (c *categoryResponseMapper) ToResponsesCategoryDeleteAt(categories []*pbcategories.CategoryResponseDeleteAt) []*response.CategoryResponseDeleteAt {
	var mappedCategories []*response.CategoryResponseDeleteAt

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategoryDelete(category))
	}

	return mappedCategories
}

func (s *categoryResponseMapper) ToResponseCashierMonthlyTotalPrice(c *pbcategories.CategoriesMonthlyTotalPriceResponse) *response.CategoriesMonthlyTotalPriceResponse {
	return &response.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToResponseCategoryMonthlyTotalPrices(c []*pbcategories.CategoriesMonthlyTotalPriceResponse) []*response.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*response.CategoriesMonthlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.ToResponseCashierMonthlyTotalPrice(Category))
	}

	return CategoryRecords
}

func (s *categoryResponseMapper) ToResponseCategoryYearlyTotalSale(c *pbcategories.CategoriesYearlyTotalPriceResponse) *response.CategoriesYearlyTotalPriceResponse {
	return &response.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToResponseCategoryYearlyTotalPrices(c []*pbcategories.CategoriesYearlyTotalPriceResponse) []*response.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*response.CategoriesYearlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.ToResponseCategoryYearlyTotalSale(Category))
	}

	return CategoryRecords
}

func (c *categoryResponseMapper) ToApiResponseCategory(pbResponse *pbcategories.ApiResponseCategory) *response.ApiResponseCategory {
	return &response.ApiResponseCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategory(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryDeleteAt(pbResponse *pbcategories.ApiResponseCategoryDeleteAt) *response.ApiResponseCategoryDeleteAt {
	return &response.ApiResponseCategoryDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryDelete(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponsesCategory(pbResponse *pbcategories.ApiResponsesCategory) *response.ApiResponsesCategory {
	return &response.ApiResponsesCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponsesCategory(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryDelete(pbResponse *pbcategories.ApiResponseCategoryDelete) *response.ApiResponseCategoryDelete {
	return &response.ApiResponseCategoryDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryAll(pbResponse *pbcategories.ApiResponseCategoryAll) *response.ApiResponseCategoryAll {
	return &response.ApiResponseCategoryAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryResponseMapper) ToApiResponsePaginationCategoryDeleteAt(pbResponse *pbcategories.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt {
	return &response.ApiResponsePaginationCategoryDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCategoryDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (c *categoryResponseMapper) ToApiResponsePaginationCategory(pbResponse *pbcategories.ApiResponsePaginationCategory) *response.ApiResponsePaginationCategory {
	return &response.ApiResponsePaginationCategory{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCategory(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryMonthlyPrice(pbResponse *pbcategories.ApiResponseCategoryMonthPrice) *response.ApiResponseCategoryMonthPrice {
	return &response.ApiResponseCategoryMonthPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryMonthlyPrices(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryYearlyPrice(pbResponse *pbcategories.ApiResponseCategoryYearPrice) *response.ApiResponseCategoryYearPrice {
	return &response.ApiResponseCategoryYearPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryYearlyPrices(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryMonthlyTotalPrice(pbResponse *pbcategories.ApiResponseCategoryMonthlyTotalPrice) *response.ApiResponseCategoryMonthlyTotalPrice {
	return &response.ApiResponseCategoryMonthlyTotalPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryMonthlyTotalPrices(pbResponse.Data),
	}
}

func (c *categoryResponseMapper) ToApiResponseCategoryYearlyTotalPrice(pbResponse *pbcategories.ApiResponseCategoryYearlyTotalPrice) *response.ApiResponseCategoryYearlyTotalPrice {
	return &response.ApiResponseCategoryYearlyTotalPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryYearlyTotalPrices(pbResponse.Data),
	}
}
