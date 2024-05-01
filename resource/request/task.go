package request

type CreateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      uint   `json:"user_id"`
}

type UpdateTask struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
