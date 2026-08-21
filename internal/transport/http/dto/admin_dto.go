package dto

type AdminCreateUserRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Nickname string  `json:"nickname"`
	City     string  `json:"city"`
	RoleIDs  []int64 `json:"roleIds"`
}

type AdminStatusRequest struct {
	Status int8 `json:"status"`
}

type AdminRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdminPermissionAssignRequest struct {
	PermissionIDs []int64 `json:"permissionIds"`
}

