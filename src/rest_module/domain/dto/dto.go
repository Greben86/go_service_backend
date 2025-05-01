package domain_dto

type ResponseHealth struct {
	Status string `json:"status"`
}

type RequestSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResponseAuth struct {
	JwtToken string `json:"token"`
}

type ResponseUser struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"email"`
}
