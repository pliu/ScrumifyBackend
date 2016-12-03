package models

import (
    "strings"
    "gopkg.in/gorp.v2"
    "TodoBackend/utils"
    "database/sql"
    "strconv"
)

type User struct {
    Id       int64    `db:"id" json:"id"`
    Username string    `db:"username" json:"username"`
    HashedPw string    `db:"hashed_pw" json:"hashed_pw"`
    Email    string    `db:"email" json:"email"`
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
    if err == sql.ErrNoRows {
        return User{}, utils.UserDoesntExist
    }
    utils.PrintErr(err, "GetUser: Failed to select user " + user_id)
    return user, err
}

func GetUsers() ([]User, error) {
    var users []User;
    _, err := Dbmap.Select(&users, "SELECT * FROM User")
    utils.PrintErr(err, "GetUsers")
    return users, err
}

func CreateUpdateUser(user User, update bool) (User, error) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "CreateUpdateUser: Failed to begin transaction")
        return User{}, err
    }
    user.Email = strings.ToLower(user.Email)

    var check User
    if err = trans.SelectOne(&check, "SELECT * FROM User WHERE email=?", user.Email); (err == nil && check.Email ==
        user.Email && update) || (err == sql.ErrNoRows && !update) {
        if update {
            _, err = trans.Update(&user)
        } else {
            err = trans.Insert(&user)
        }
        if err == nil {
            return user, trans.Commit()
        } else {
            trans.Rollback()
            utils.PrintErr(err, "CreateUpdateUser: Failed to insert/update user " + strconv.FormatInt(user.Id, 10))
            return User{}, err
        }
    } else if err != nil {
        utils.PrintErr(err, "CreateUpdateUser: Failed to select email " + user.Email)
        return User{}, err
    } else {
        return User{}, utils.EmailExists
    }
}

func DeleteUser(user User) error {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "DeleteUser: Failed to begin transaction")
        return err
    }

    var mappings []EpicUserMap
    if _, err = trans.Select(&mappings, "SELECT * FROM EpicUserMap WHERE user_id=?", user.Id); err == nil {
        if _, err = trans.Delete(&user); err == nil {
            if _, err = trans.Exec("DELETE FROM EpicUserMap WHERE user_id=?", user.Id); err == nil {
                err = trans.Commit()
            } else {
                trans.Rollback()
                utils.PrintErr(err, "DeleteUser: Failed to delete mappings for user " + strconv.FormatInt(user.Id, 10))
                return err
            }
        } else {
            trans.Rollback()
            utils.PrintErr(err, "DeleteUser: Failed to delete user " + strconv.FormatInt(user.Id, 10))
            return err
        }
        for _, mapping := range mappings {
            go removeUnownedEpic(mapping.EpicId)
        }
        return err
    } else {
        utils.PrintErr(err, "DeleteUser: Failed to select user " + strconv.FormatInt(user.Id, 10))
        return err
    }
}

func (user User)IsValid() bool {
    if user.Username != "" && user.HashedPw != "" && user.Email != "" {
        return true
    } else {
        return false
    }
}
