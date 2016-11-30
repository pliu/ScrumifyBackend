package models

import (
	"strings"
	"gopkg.in/gorp.v2"
)

type User struct {
	Id       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	HashedPw string `db:"hashedpw" json:"hashedpw"`
	Email    string `db:"email" json:"email"`
}

func SetUserProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.ColMap("Username").SetNotNull(true)
	table.ColMap("HashedPw").SetNotNull(true)
	table.ColMap("Email").SetUnique(true).SetNotNull(true)
}

func GetUser(user_id string) (User, error) {
	var user User
	err := Dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", user_id)
	return user, err
}

func GetUsers() ([]User, error) {
	var users []User;
	_, err := Dbmap.Select(&users, "SELECT * FROM User")
	return users, err
}

func CreateUser(user User) (User, error) {
	user.Email = strings.ToLower(user.Email)
	err := Dbmap.Insert(&user)
	return ScrubUser(user), err
}

func UpdateUser(user User) (User, error) {
	user.Email = strings.ToLower(user.Email)
	_, err := Dbmap.Update(&user)
	return ScrubUser(user), err
}

func DeleteUser(user User) error {
	trans, err := Dbmap.Begin()
	if err != nil {
		return err
	}

	var mappings []EpicUserMap
	if _, err = trans.Select(mappings, "SELECT * FROM EpicUserMap WHERE userid=?", user.Id); err == nil {
		if _, err = trans.Delete(&user); err == nil {
			if _, err = trans.Exec("DELETE FROM EpicUserMap WHERE userid=?", user.Id); err == nil {
				err = trans.Commit()
			} else {
				trans.Rollback()
				return err
			}
		} else {
			trans.Rollback()
			return err
		}
		/*for _, mapping := range mappings {

		}*/
		return err
	} else {
		return err
	}
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
