package viewmodels

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserView struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	UserView
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	UserView
}
