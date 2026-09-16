package schemas

type CreateBookRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AuthorNames []string `json:"authorNames"`
}
