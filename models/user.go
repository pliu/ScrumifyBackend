package models

type User struct {
	Id       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	HashedPW string `db:"hashedpw" json:"hashedpw"`
	Email    string `db:"email" json:"email"`
}

func GetUser() {

}

func CreateUser() {

}

func UpdateUser() {

}

func DeleteUser() {

}
