package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username     string `json:"username"`
	UserId       string `json:"userId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Permission   string `json:"permission,omitempty"`
}

type CreateUserRequest struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
	Password   string `json:"password"`
}

type DeleteUserRequest struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
	Password   string `json:"password"`
}
