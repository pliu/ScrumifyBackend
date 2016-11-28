package models

import (
	"strings"
	"gopkg.in/gorp.v2"
	"errors"
)

type User struct {
	Id       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	HashedPw string `db:"hashedpw" json:"hashedpw"`
	Email    string `db:"email" json:"email"`
}

func SetUserProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.AddIndex("EmailIndex", "Hash", []string{"Email"})
	table.ColMap("Username").SetNotNull(true)
	table.ColMap("HashedPw").SetNotNull(true)
	table.ColMap("Email").SetUnique(true).SetNotNull(true)
}

func GetUser(id string) (User, error) {
	var user User
	err := Dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)
	return user, err
}

func GetUsers() ([]User, error) {
	var users []User;
	_, err := Dbmap.Select(&users, "SELECT * FROM User")
	return users, err
}

func CreateUser(user User) (User, error) {
	user.Email = strings.ToLower(user.Email)
	if insert, _ := Dbmap.Exec(`INSERT INTO User (username, hashedpw, email) VALUES (?, ?, ?)`, user.Username,
			user.HashedPw, user.Email); insert != nil {
		user_id, nerr := insert.LastInsertId()
		if nerr == nil {
			user.Id = user_id
		}
		return ScrubUser(user), nerr
	} else {
		return ScrubUser(user), errors.New("Failed to insert user into database")
	}
}

func UpdateUser(user User) (User, error) {
	user.Email = strings.ToLower(user.Email)
	_, err := Dbmap.Update(&user)
	return ScrubUser(user), err
}

func DeleteUser(user User) error {
	_, err := Dbmap.Delete(&user)
	return err
}

func GetUserByEmail(email string) (User, error) {
	email = strings.ToLower(email)
	var user User
	err := Dbmap.SelectOne(&user, "SELECT id FROM User WHERE email=?", email)
	return user, err
}

func (user User)IsValid() bool {
	if user.Username != "" && user.HashedPw != "" && user.Email != "" {
		return true
	} else {
		return false
	}
}

func ScrubUser(user User) User {
	return User{
		Id:       user.Id,
		Username: user.Username,
		HashedPw: "",
		Email:    user.Email,
	}
}
