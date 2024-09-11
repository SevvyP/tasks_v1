package model

type Task struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Body      string  `json:"body"`
	Completed bool    `json:"completed"`
	Parent    *string `json:"parent"`
	Reminder  *string `json:"reminder"`
}
