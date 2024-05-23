package web

type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"ConfirmPassword"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EditReq struct {
	NickName        string
	BirthDate       string
	PersonalProfile string
}

type SendSmsCodeReq struct {
	Phone string `json:"phone"`
}

type VerifySmsCodeReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
