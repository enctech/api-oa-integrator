package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username     string   `json:"username"`
	UserId       string   `json:"userId"`
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	Permissions  []string `json:"permissions,omitempty"`
}

type UsersResponse struct {
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	UserId      string   `json:"userId"`
	Permissions []string `json:"permissions,omitempty"`
}

type CreateUserRequest struct {
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	Password    string   `json:"password"`
}

type DeleteUserRequest struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
	Password   string `json:"password"`
}
