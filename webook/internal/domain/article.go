package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type ArticleStatus uint8

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

const (
	// ArticleStatusUnknown 这是一个未知状态
	ArticleStatusUnknown     = iota
	ArticleStatusUnpublished = iota
	ArticleStatusPublished
	ArticleStatusPrivate
)

type Author struct {
	Id   int64
	Name string
}
