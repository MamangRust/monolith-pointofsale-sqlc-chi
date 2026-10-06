package response_api

import (
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type userResponseMapper struct {
}

func NewUserResponseMapper() *userResponseMapper {
	return &userResponseMapper{}
}

func (u *userResponseMapper) ToResponseUser(user *pbusers.UserResponse) *response.UserResponse {
	return &response.UserResponse{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userResponseMapper) ToResponsesUser(users []*pbusers.UserResponse) []*response.UserResponse {
	var mappedUsers []*response.UserResponse

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.ToResponseUser(user))
	}

	return mappedUsers
}

func (u *userResponseMapper) ToResponseUserDeleteAt(user *pbusers.UserResponseDeleteAt) *response.UserResponseDeleteAt {
	var deletedAt string
	if user.DeletedAt != nil {
		deletedAt = user.DeletedAt.Value
	}

	return &response.UserResponseDeleteAt{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (u *userResponseMapper) ToResponsesUserDeleteAt(users []*pbusers.UserResponseDeleteAt) []*response.UserResponseDeleteAt {
	var mappedUsers []*response.UserResponseDeleteAt

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.ToResponseUserDeleteAt(user))
	}

	return mappedUsers
}

func (u *userResponseMapper) ToApiResponseUser(pbResponse *pbusers.ApiResponseUser) *response.ApiResponseUser {
	return &response.ApiResponseUser{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUser(pbResponse.Data),
	}
}

func (u *userResponseMapper) ToApiResponseUserDeleteAt(pbResponse *pbusers.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt {
	return &response.ApiResponseUserDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUserDeleteAt(pbResponse.Data),
	}
}

func (u *userResponseMapper) ToApiResponsesUser(pbResponse *pbusers.ApiResponsesUser) *response.ApiResponsesUser {
	return &response.ApiResponsesUser{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponsesUser(pbResponse.Data),
	}
}

func (u *userResponseMapper) ToApiResponseUserDelete(pbResponse *pbusers.ApiResponseUserDelete) *response.ApiResponseUserDelete {
	return &response.ApiResponseUserDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (u *userResponseMapper) ToApiResponseUserAll(pbResponse *pbusers.ApiResponseUserAll) *response.ApiResponseUserAll {
	return &response.ApiResponseUserAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (u *userResponseMapper) ToApiResponsePaginationUserDeleteAt(pbResponse *pbusers.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt {
	return &response.ApiResponsePaginationUserDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       u.ToResponsesUserDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (u *userResponseMapper) ToApiResponsePaginationUser(pbResponse *pbusers.ApiResponsePaginationUser) *response.ApiResponsePaginationUser {
	return &response.ApiResponsePaginationUser{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       u.ToResponsesUser(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
