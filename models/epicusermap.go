package models

import (
    "gopkg.in/gorp.v2"
    "database/sql"
    "ScrumifyBackend/utils"
    "strconv"
)

// Which users are part of which epics
type EpicUserMap struct {
    UserId int64 `db:"user_id" json:"user_id"`
    EpicId int64 `db:"epic_id" json:"epic_id"`
}

func SetEpicUserMapProperties(table *gorp.TableMap) {
    table.SetKeys(false, "UserId", "EpicId")

    // InnoDB does not have Hash indices
    table.AddIndex("MapEpicIdIndex", "Btree", []string{"EpicId"})
}

func AddEpicUserMap(email string, epic_id string) (User, error) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "AddEpicUserMap: Failed to begin transaction")
        return User{}, err
    }

    var userToAdd User
    if err = trans.SelectOne(&userToAdd, "SELECT * FROM User WHERE email=?", email); err == nil {
        if _, err := strconv.ParseInt(epic_id, 10, 64); err == nil {
            if _, err = trans.Exec(`INSERT INTO EpicUserMap (user_id, epic_id) VALUES (?, ?)`, userToAdd.Id, epic_id);
                err == nil {
                return userToAdd, trans.Commit()
            } else {
                trans.Rollback()
                utils.PrintErr(err, "AddEpicUserMap: Failed to insert mapping user_id: " +
                    strconv.FormatInt(userToAdd.Id, 10) + " epic_id: " + epic_id)
                return User{}, err  // TODO: Differentiate between collision errors and just DB failure
            }
        } else {
            trans.Rollback()
            return User{}, utils.UserDoesntExist
        }
    } else if err == sql.ErrNoRows {
        trans.Rollback()
        return User{}, utils.EmailDoesntExist
    } else {
        trans.Rollback()
        utils.PrintErr(err, "AddEpicUserMap: Failed to select email " + email)
        return User{}, err
    }
}

func EpicOwnedByUser(user_id string, epic_id string) (EpicUserMap, error) {
    var epicUserMap EpicUserMap
    err := Dbmap.SelectOne(&epicUserMap, "SELECT * FROM EpicUserMap WHERE user_id=? AND epic_id=?", user_id, epic_id)
    if err == sql.ErrNoRows {
        return EpicUserMap{}, utils.MappingDoesntExist
    }
    utils.PrintErr(err, "EpicOwnedByUser: Failed to select mapping user_id: " + user_id + " epic_id: " + epic_id)
    return epicUserMap, err
}
