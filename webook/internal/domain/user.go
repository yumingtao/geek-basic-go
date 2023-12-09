package domain

type User struct {
	Id              int64
	Email           string
	Password        string `json:"-"`
	NickName        string
	BirthDate       string
	PersonalProfile string
	Phone           string
}

// 按照DDD的原则，User password和email的校验应该放在这里
/*func (u User) isEmailValid() bool {
	return u.Email
}*/
