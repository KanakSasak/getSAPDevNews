package model

type News struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Desc     string `json:"desc"`
	Link     string `json:"link"`
}
