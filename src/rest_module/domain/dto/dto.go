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

type ResponseTransfer struct {
	AccountFrom int64   `json:"account_from"`
	AccountTo   int64   `json:"account_to"`
	SumValue    float64 `json:"sum_value"`
}
