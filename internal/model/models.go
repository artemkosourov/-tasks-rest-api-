package model

type User struct {
	ID       int    `json:"id" example:"1"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"-"`
	Role     string `json:"role" example:"user"`
}

type Task struct {
	ID        int    `json:"id" example:"1"`
	Title     string `json:"title" example:"Купить продукты"`
	UserID    int    `json:"user_id" example:"1"`
	Completed bool   `json:"completed" example:"false"`
}
