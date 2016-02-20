package model

type User struct {
	Id int64 `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname string `db:"lastname" json:"lastname"`
}
