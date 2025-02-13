package req

type AdminQuery struct {
	Email       *string `query:"email"`
	PartialName *string `query:"partial_name"`
	Pagination
}

type CreateAdminBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
