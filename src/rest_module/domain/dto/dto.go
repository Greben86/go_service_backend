package domain_dto

type RequestSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseAuth struct {
	JwtToken string `json:"token"`
}

type ResponseUser struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
}
