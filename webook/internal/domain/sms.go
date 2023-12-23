package domain

type Sms struct {
	Id          int64
	Tpl         string
	Args        []string
	Numbers     []string
	RetryMaxCnt int
}
