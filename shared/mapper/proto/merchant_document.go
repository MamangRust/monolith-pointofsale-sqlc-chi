package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantDocumentProtoMapper struct{}

func NewMerchantDocumentProtoMapper() *merchantDocumentProtoMapper {
	return &merchantDocumentProtoMapper{}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pbmerchants_document.ApiResponseMerchantDocument {
	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocument(doc),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchants_document.ApiResponsesMerchantDocument {
	return &pbmerchants_document.ApiResponsesMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocuments(docs),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchants_document.ApiResponsePaginationMerchantDocument {
	return &pbmerchants_document.ApiResponsePaginationMerchantDocument{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocuments(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pbmerchants_document.ApiResponsePaginationMerchantDocumentAt {
	return &pbmerchants_document.ApiResponsePaginationMerchantDocumentAt{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocumentsDeleteAt(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentDelete(status string, message string) *pbmerchants_document.ApiResponseMerchantDocumentDelete {
	return &pbmerchants_document.ApiResponseMerchantDocumentDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentAll(status string, message string) *pbmerchants_document.ApiResponseMerchantDocumentAll {
	return &pbmerchants_document.ApiResponseMerchantDocumentAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocument(doc *response.MerchantDocumentResponse) *pbmerchants_document.MerchantDocument {
	return &pbmerchants_document.MerchantDocument{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocuments(docs []*response.MerchantDocumentResponse) []*pbmerchants_document.MerchantDocument {
	var res []*pbmerchants_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocument(doc))
	}
	return res
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentDeleteAt(doc *response.MerchantDocumentResponseDeleteAt) *pbmerchants_document.MerchantDocumentDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if doc.DeletedAt != nil {
		deletedAt = wrapperspb.String(*doc.DeletedAt)
	}

	return &pbmerchants_document.MerchantDocumentDeleteAt{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentsDeleteAt(docs []*response.MerchantDocumentResponseDeleteAt) []*pbmerchants_document.MerchantDocumentDeleteAt {
	var res []*pbmerchants_document.MerchantDocumentDeleteAt
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocumentDeleteAt(doc))
	}
	return res
}
