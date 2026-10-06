package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
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

func mapGetOrderItemsRowToProto(orderItem *models.OrderItemRow) *pborder_items.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !orderItem.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&orderItem.CreatedAt)
	}
	if !orderItem.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&orderItem.UpdatedAt)
	}
	return &pborder_items.OrderItemResponse{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesOrderItem(orderItems []*models.OrderItemRow) []*pborder_items.OrderItemResponse {
	var mappedOrderItems []*pborder_items.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapGetOrderItemsActiveRowToProto(orderItem *models.OrderItemActiveRow) *pborder_items.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if !orderItem.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&orderItem.CreatedAt)
	}
	if !orderItem.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&orderItem.UpdatedAt)
	}
	if orderItem.DeletedAt != nil {
		deletedAt = wrapperspb.String(orderItem.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pborder_items.OrderItemResponseDeleteAt{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapResponsesOrderItemActive(orderItems []*models.OrderItemActiveRow) []*pborder_items.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pborder_items.OrderItemResponseDeleteAt
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsActiveRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapGetOrderItemsTrashedRowToProto(orderItem *models.OrderItemTrashedRow) *pborder_items.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if !orderItem.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&orderItem.CreatedAt)
	}
	if !orderItem.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&orderItem.UpdatedAt)
	}
	if orderItem.DeletedAt != nil {
		deletedAt = wrapperspb.String(orderItem.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pborder_items.OrderItemResponseDeleteAt{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapResponsesOrderItemTrashed(orderItems []*models.OrderItemTrashedRow) []*pborder_items.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pborder_items.OrderItemResponseDeleteAt
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsTrashedRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapOrderItemToProto(orderItem *models.OrderItem) *pborder_items.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !orderItem.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&orderItem.CreatedAt)
	}
	if !orderItem.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&orderItem.UpdatedAt)
	}
	return &pborder_items.OrderItemResponse{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesOrderItemFromModel(orderItems []*models.OrderItem) []*pborder_items.OrderItemResponse {
	var mappedOrderItems []*pborder_items.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapOrderItemToProto(orderItem))
	}
	return mappedOrderItems
}
