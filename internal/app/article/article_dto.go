package article

type ListArticleDTO struct {
	Articles []Article `json:"items"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}
