package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type productProtoMapper struct{}

func NewProductProtoMapper() *productProtoMapper {
	return &productProtoMapper{}
}

func (p *productProtoMapper) ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproducts.ApiResponseProduct {
	return &pbproducts.ApiResponseProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pbproducts.ApiResponsesProduct {
	return &pbproducts.ApiResponsesProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponsesProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproducts.ApiResponseProductDeleteAt {
	return &pbproducts.ApiResponseProductDeleteAt{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProductDeleteAt(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDelete(status string, message string) *pbproducts.ApiResponseProductDelete {
	return &pbproducts.ApiResponseProductDelete{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponseProductAll(status string, message string) *pbproducts.ApiResponseProductAll {
	return &pbproducts.ApiResponseProductAll{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pbproducts.ApiResponsePaginationProductDeleteAt {
	return &pbproducts.ApiResponsePaginationProductDeleteAt{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProductDeleteAt(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pbproducts.ApiResponsePaginationProduct {
	return &pbproducts.ApiResponsePaginationProduct{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProduct(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) mapResponseProduct(product *response.ProductResponse) *pbproducts.ProductResponse {
	return &pbproducts.ProductResponse{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productProtoMapper) mapResponsesProduct(products []*response.ProductResponse) []*pbproducts.ProductResponse {
	var mappedProducts []*pbproducts.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProduct(product))
	}

	return mappedProducts
}

func (p *productProtoMapper) mapResponseProductDeleteAt(product *response.ProductResponseDeleteAt) *pbproducts.ProductResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if product.DeleteAt != nil {
		deletedAt = wrapperspb.String(*product.DeleteAt)
	}

	return &pbproducts.ProductResponseDeleteAt{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (p *productProtoMapper) mapResponsesProductDeleteAt(products []*response.ProductResponseDeleteAt) []*pbproducts.ProductResponseDeleteAt {
	var mappedProducts []*pbproducts.ProductResponseDeleteAt

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProductDeleteAt(product))
	}

	return mappedProducts
}
