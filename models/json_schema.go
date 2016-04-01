package models

type ModuleIn struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	DueDate      string  `json:"duedate"`
	Stage        int64   `json:"stage"`
	Dependencies []int64 `json:"dependencies"`
}

type EmailIn struct {
	Email string `json:"email"`
}
