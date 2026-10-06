package handler

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
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

func mapResponseCategory(category *models.Category) *pbcategories.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&category.UpdatedAt)
	}
	return &pbcategories.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesCategoryModels(categories []*models.Category) []*pbcategories.CategoryResponse {
	mappedCategories := make([]*pbcategories.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseCategory(category))
	}
	return mappedCategories
}

func mapResponseGetCategory(category *models.CategoryRow) *pbcategories.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&category.UpdatedAt)
	}
	return &pbcategories.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesCategory(categories []*models.CategoryRow) []*pbcategories.CategoryResponse {
	var mappedCategories []*pbcategories.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategory(category))
	}
	return mappedCategories
}

func mapResponseCategoryDeleteAt(category *models.Category) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = convert.TimeToWrappers(&category.DeletedAt.Time)
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponseGetCategoryActive(category *models.CategoryActiveRow) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(category.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryActive(categories []*models.CategoryActiveRow) []*pbcategories.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategories.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryActive(category))
	}
	return mappedCategories
}

func mapResponseGetCategoryTrashed(category *models.CategoryTrashedRow) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(category.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryTrashed(categories []*models.CategoryTrashedRow) []*pbcategories.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategories.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryTrashed(category))
	}
	return mappedCategories
}

func mapResponseCategoryMonthlyPrice(category *models.CategoryMonthlyRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPrices(c []*models.CategoryMonthlyRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceById(category *models.CategoryMonthlyByIdRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesById(c []*models.CategoryMonthlyByIdRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceByMerchant(category *models.CategoryMonthlyByMerchantRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesByMerchant(c []*models.CategoryMonthlyByMerchantRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPrice(category *models.CategoryYearlyRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPrices(c []*models.CategoryYearlyRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceById(category *models.CategoryYearlyByIdRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesById(c []*models.CategoryYearlyByIdRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceByMerchant(category *models.CategoryYearlyByMerchantRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesByMerchant(c []*models.CategoryYearlyByMerchantRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCashierMonthlyTotalPrice(c *models.CategoryMonthlyTotalPriceRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPrices(c []*models.CategoryMonthlyTotalPriceRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPrice(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceById(c *models.CategoryMonthlyTotalPriceByIdRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesById(c []*models.CategoryMonthlyTotalPriceByIdRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceById(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceByMerchant(c *models.CategoryMonthlyTotalPriceByMerchantRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesByMerchant(c []*models.CategoryMonthlyTotalPriceByMerchantRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceByMerchant(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSale(c *models.CategoryYearlyTotalPriceRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPrices(c []*models.CategoryYearlyTotalPriceRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSale(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleById(c *models.CategoryYearlyTotalPriceByIdRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesById(c []*models.CategoryYearlyTotalPriceByIdRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleById(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleByMerchant(c *models.CategoryYearlyTotalPriceByMerchantRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesByMerchant(c []*models.CategoryYearlyTotalPriceByMerchantRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleByMerchant(Category))
	}
	return CategoryRecords
}
