package model

type Category struct {
	Id        int64
	ParentId  int64
	Title     string
	MetaTitle string
	Slug      string
	Content   string
	Posts     *[]Post
}
