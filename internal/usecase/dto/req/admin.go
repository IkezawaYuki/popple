package req

type AdminQuery struct {
	PartialName string `query:"partial_name" binding:"required"`
}

type CreateAdminBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
