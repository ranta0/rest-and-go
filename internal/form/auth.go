package form

type Auth struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type AuthRefresh struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
}
