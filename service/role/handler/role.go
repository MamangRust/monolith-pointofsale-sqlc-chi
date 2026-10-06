package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
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

func mapResponseRole(role *models.Role) *pbroles.RoleResponse {
	if role == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(&role.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(&role.UpdatedAt)
	return &pbroles.RoleResponse{
		Id:        role.RoleID,
		Name:      role.RoleName,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapUserRoleToProto(userRole *models.UserRole) *pbroles.RoleUserResponse {
	if userRole == nil {
		return nil
	}
	return &pbroles.RoleUserResponse{
		UserRoleId: userRole.UserRoleID,
		UserId:     userRole.UserID,
		RoleId:     userRole.RoleID,
	}
}

func mapResponsesRole(roles []*models.RoleRow) []*pbroles.RoleResponse {
	var responseRoles []*pbroles.RoleResponse
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponse{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(&role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(&role.UpdatedAt),
		})
	}
	return responseRoles
}

func mapResponsesRoleFromDB(roles []*models.Role) []*pbroles.RoleResponse {
	var responseRoles []*pbroles.RoleResponse
	for _, role := range roles {
		responseRoles = append(responseRoles, mapResponseRole(role))
	}
	return responseRoles
}

func formatOptionalTime(t *time.Time) string {
	return convert.FormatTimePtr(t)
}

func mapResponsesRoleFromActive(roles []*models.RoleActiveRow) []*pbroles.RoleResponseDeleteAt {
	var responseRoles []*pbroles.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(&role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(&role.UpdatedAt),
			DeletedAt: formatOptionalTime(role.DeletedAt),
		})
	}
	return responseRoles
}

func mapResponsesRoleFromTrashed(roles []*models.RoleTrashedRow) []*pbroles.RoleResponseDeleteAt {
	var responseRoles []*pbroles.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(&role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(&role.UpdatedAt),
			DeletedAt: formatOptionalTime(role.DeletedAt),
		})
	}
	return responseRoles
}
