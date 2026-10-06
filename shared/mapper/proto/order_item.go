package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderItemProtoMapper struct{}

func NewOrderItemProtoMapper() *orderItemProtoMapper {
	return &orderItemProtoMapper{}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pborder_items.ApiResponseOrderItem {
	return &pborder_items.ApiResponseOrderItem{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderItem(pbResponse),
	}
}

func (o *orderItemProtoMapper) ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pborder_items.ApiResponsesOrderItem {
	return &pborder_items.ApiResponsesOrderItem{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderItem(pbResponse),
	}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItemDelete(status string, message string) *pborder_items.ApiResponseOrderItemDelete {
	return &pborder_items.ApiResponseOrderItemDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItemAll(status string, message string) *pborder_items.ApiResponseOrderItemAll {
	return &pborder_items.ApiResponseOrderItemAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderItemProtoMapper) ToProtoResponsePaginationOrderItemDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pborder_items.ApiResponsePaginationOrderItemDeleteAt {
	return &pborder_items.ApiResponsePaginationOrderItemDeleteAt{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderItemDeleteAt(orderItems),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderItemProtoMapper) ToProtoResponsePaginationOrderItem(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pborder_items.ApiResponsePaginationOrderItem {
	return &pborder_items.ApiResponsePaginationOrderItem{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderItem(orderItems),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderItemProtoMapper) mapResponseOrderItem(orderItem *response.OrderItemResponse) *pborder_items.OrderItemResponse {
	return &pborder_items.OrderItemResponse{
		Id:        int32(orderItem.ID),
		OrderId:   int32(orderItem.OrderID),
		ProductId: int32(orderItem.ProductID),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}
}

func (o *orderItemProtoMapper) mapResponsesOrderItem(orderItems []*response.OrderItemResponse) []*pborder_items.OrderItemResponse {
	var mappedOrderItems []*pborder_items.OrderItemResponse

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.mapResponseOrderItem(orderItem))
	}

	return mappedOrderItems
}

func (o *orderItemProtoMapper) mapResponseOrderItemDelete(orderItem *response.OrderItemResponseDeleteAt) *pborder_items.OrderItemResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if orderItem.DeleteAt != nil {
		deletedAt = wrapperspb.String(*orderItem.DeleteAt)
	}

	return &pborder_items.OrderItemResponseDeleteAt{
		Id:        int32(orderItem.ID),
		OrderId:   int32(orderItem.OrderID),
		ProductId: int32(orderItem.ProductID),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (o *orderItemProtoMapper) mapResponsesOrderItemDeleteAt(orderItems []*response.OrderItemResponseDeleteAt) []*pborder_items.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pborder_items.OrderItemResponseDeleteAt

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.mapResponseOrderItemDelete(orderItem))
	}

	return mappedOrderItems
}
