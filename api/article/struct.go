package article

type CreateArticleParams struct {
	Title         string   `json:"title"`
	Author        int      `json:"author"`
	Content       string   `json:"content"`
	Private       bool     `json:"private"`
	ArticleLabels []string `json:"labels"`
	ArticleTypes  []string `json:"types"`
}
