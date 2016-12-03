package models

type RestStory struct {
    Id           int64   `json:"id"`
    Name         string  `json:"name"`
    DueDate      string  `json:"duedate"`
    Stage        int64   `json:"stage"`
    EpicId       int64    `json:"epicid"`
    Dependencies []int64    `json:"dependencies"`
}

type RestEmail struct {
    Email string `json:"email"`
}
