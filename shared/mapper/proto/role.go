package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type roleProtoMapper struct {
}

func NewRoleProtoMapper() *roleProtoMapper {
	return &roleProtoMapper{}
}

func (s *roleProtoMapper) ToProtoResponseRoleAll(status string, message string) *pbroles.ApiResponseRoleAll {
	return &pbroles.ApiResponseRoleAll{
		Status:  status,
		Message: message,
	}
}

func (s *roleProtoMapper) ToProtoResponseRoleDelete(status string, message string) *pbroles.ApiResponseRoleDelete {
	return &pbroles.ApiResponseRoleDelete{
		Status:  status,
		Message: message,
	}
}

func (s *roleProtoMapper) ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbroles.ApiResponseRole {
	return &pbroles.ApiResponseRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponseRole(pbResponse),
	}
}

func (s *roleProtoMapper) ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbroles.ApiResponsesRole {
	return &pbroles.ApiResponsesRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesRole(pbResponse),
	}
}

func (s *roleProtoMapper) ToProtoResponsePaginationRole(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pbroles.ApiResponsePaginationRole {
	return &pbroles.ApiResponsePaginationRole{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesRole(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *roleProtoMapper) ToProtoResponsePaginationRoleDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbroles.ApiResponsePaginationRoleDeleteAt {
	return &pbroles.ApiResponsePaginationRoleDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *roleProtoMapper) mapResponseRole(role *response.RoleResponse) *pbroles.RoleResponse {
	return &pbroles.RoleResponse{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleProtoMapper) mapResponsesRole(roles []*response.RoleResponse) []*pbroles.RoleResponse {
	var responseRoles []*pbroles.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleProtoMapper) mapResponseRoleDeleteAt(role *response.RoleResponseDeleteAt) *pbroles.RoleResponseDeleteAt {
	return &pbroles.RoleResponseDeleteAt{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}

func (s *roleProtoMapper) mapResponsesRoleDeleteAt(roles []*response.RoleResponseDeleteAt) []*pbroles.RoleResponseDeleteAt {
	var responseRoles []*pbroles.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}
