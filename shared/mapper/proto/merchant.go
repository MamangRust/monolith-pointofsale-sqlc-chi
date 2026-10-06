package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type merchantProtoMapper struct{}

func NewMerchantProtoMaper() *merchantProtoMapper {
	return &merchantProtoMapper{}
}

func (m *merchantProtoMapper) ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pbmerchants.ApiResponseMerchant {
	return &pbmerchants.ApiResponseMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pbmerchants.ApiResponsesMerchant {
	return &pbmerchants.ApiResponsesMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pbmerchants.ApiResponseMerchantDeleteAt {
	return &pbmerchants.ApiResponseMerchantDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDeleteAt(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDelete(status string, message string) *pbmerchants.ApiResponseMerchantDelete {
	return &pbmerchants.ApiResponseMerchantDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantAll(status string, message string) *pbmerchants.ApiResponseMerchantAll {
	return &pbmerchants.ApiResponseMerchantAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pbmerchants.ApiResponsePaginationMerchantDeleteAt {
	return &pbmerchants.ApiResponsePaginationMerchantDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantDeleteAt(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pbmerchants.ApiResponsePaginationMerchant {
	return &pbmerchants.ApiResponsePaginationMerchant{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchant(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) mapResponseMerchant(merchant *response.MerchantResponse) *pbmerchants.MerchantResponse {
	return &pbmerchants.MerchantResponse{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
	}
}

func (m *merchantProtoMapper) mapResponsesMerchant(merchants []*response.MerchantResponse) []*pbmerchants.MerchantResponse {
	var mappedMerchants []*pbmerchants.MerchantResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchant(merchant))
	}

	return mappedMerchants
}

func (m *merchantProtoMapper) mapResponseMerchantDeleteAt(merchant *response.MerchantResponseDeleteAt) *pbmerchants.MerchantResponseDeleteAt {
	return &pbmerchants.MerchantResponseDeleteAt{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		DeletedAt:    merchant.DeletedAt,
	}
}

func (m *merchantProtoMapper) mapResponsesMerchantDeleteAt(merchants []*response.MerchantResponseDeleteAt) []*pbmerchants.MerchantResponseDeleteAt {
	var mappedMerchants []*pbmerchants.MerchantResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantDeleteAt(merchant))
	}

	return mappedMerchants
}
