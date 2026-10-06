package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"

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

func formatOptionalTime(t *time.Time) string {
	return convert.FormatTimePtr(t)
}

func mapUserToProto(user *models.User) *pbusers.UserResponse {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	return &pbusers.UserResponse{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

// mapUserWithPasswordToProto exposes the password hash. It is only used by the
// auth-oriented RPCs and must never back a public listing endpoint.
func mapUserWithPasswordToProto(user *models.User) *pbusers.UserResponseWithPassword {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	return &pbusers.UserResponseWithPassword{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapGetUsersRowToProto(user *models.UserRow) *pbusers.UserResponse {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	return &pbusers.UserResponse{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapGetUsersRowsToProto(users []*models.UserRow) []*pbusers.UserResponse {
	var responseUsers []*pbusers.UserResponse
	for _, u := range users {
		responseUsers = append(responseUsers, mapGetUsersRowToProto(u))
	}
	return responseUsers
}

func mapUserDeleteAtToProto(user *models.User) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	var deletedAt *wrapperspb.StringValue
	if user.DeletedAt.Valid {
		deletedAt = convert.TimeToWrappers(&user.DeletedAt.Time)
	}
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapActiveUserToProto(user *models.UserActiveRow) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	deletedAtVal := convert.TimeToWrappers(user.DeletedAt)
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAtVal,
	}
}

func mapActiveUsersToProto(users []*models.UserActiveRow) []*pbusers.UserResponseDeleteAt {
	var responseUsers []*pbusers.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapActiveUserToProto(u))
	}
	return responseUsers
}

func mapTrashedUserToProto(user *models.UserTrashedRow) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&user.UpdatedAt)
	deletedAtVal := convert.TimeToWrappers(user.DeletedAt)
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAtVal,
	}
}

func mapTrashedUsersToProto(users []*models.UserTrashedRow) []*pbusers.UserResponseDeleteAt {
	var responseUsers []*pbusers.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapTrashedUserToProto(u))
	}
	return responseUsers
}
