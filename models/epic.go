package models

import (
    "gopkg.in/gorp.v2"
    "ScrumifyBackend/utils"
    "strconv"
    "time"
)

type Epic struct {
    Id        int64     `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Members   []User    `db:"-" json:"members"`
    Stories   []Story   `db:"-" json:"stories"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func SetEpicProperties(table *gorp.TableMap) {
    table.SetKeys(true, "Id")

    // InnoDB does not have Hash indices
    table.AddIndex("EpicCreatedAtIndex", "Btree", []string{"created_at"})
    table.AddIndex("EpicUpdatedAtIndex", "Btree", []string{"updated_at"})
    table.ColMap("Name").SetNotNull(true)
    table.ColMap("CreatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP")
    table.ColMap("UpdatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP ON UPDATE " +
        "CURRENT_TIMESTAMP")
}

func GetEpics(user_id string) ([]Epic, error) {
    var epics []Epic
    _, err := Dbmap.Select(&epics, "SELECT * FROM Epic WHERE id IN (SELECT epic_id FROM EpicUserMap WHERE user_id=?)",
        user_id)
    utils.PrintErr(err, "GetEpics: Failed to select epics for user " + user_id)
    return epics, err
}

func CreateEpic(user_id string, epic Epic) (Epic, error) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "CreateEpic: Failed to begin transaction")
        return Epic{}, err
    }

    if err = trans.Insert(&epic); err == nil {
        if _, err = trans.Exec("INSERT INTO EpicUserMap (user_id, epic_id) VALUES (?, ?)", user_id, epic.Id);
            err == nil {
            var check Epic
            if err = trans.SelectOne(&check, "SELECT * FROM Epic WHERE id=?", epic.Id); err == nil {
                return check, trans.Commit()
            } else {
                return epic, trans.Commit()
            }
        } else {
            utils.PrintErr(err, "CreateEpic: Failed to insert mapping user_id: " + user_id + " epic_id: " +
                strconv.FormatInt(epic.Id, 10))
        }
    } else {
        utils.PrintErr(err, "CreateEpic: Failed to insert epic " + strconv.FormatInt(epic.Id, 10))
    }
    trans.Rollback()
    return Epic{}, err
}

func UpdateEpic(epic Epic) (Epic, error) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "UpdateEpic: Failed to begin transaction")
        return Epic{}, err
    }

    if _, err = trans.Update(&epic); err == nil {
        var check Epic
        if err = trans.SelectOne(&check, "SELECT * FROM Epic WHERE id=?", epic.Id); err == nil {
            return check, trans.Commit()
        } else {
            return epic, trans.Commit()
        }
    } else {
        trans.Rollback()
        utils.PrintErr(err, "UpdateEpic: Failed to update epic " + strconv.FormatInt(epic.Id, 10))
        return Epic{}, err
    }
}

func DeleteEpic(mapping EpicUserMap) error {
    _, err := Dbmap.Delete(&mapping)
    utils.PrintErr(err, "DeleteEpic: Failed to delete mapping user_id: " + strconv.FormatInt(mapping.UserId, 10) +
        " epic_id: " + strconv.FormatInt(mapping.EpicId, 10))
    if err == nil {
        go removeUnownedEpic(mapping.EpicId)
    }
    return err
}

func (epic Epic)IsValid() bool {
    return epic.Name != ""
}

// Called as a goroutine
func removeUnownedEpic(epic_id int64) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "removeUnownedEpic: Failed to begin transaction")
        return
    }

    var count int64
    if count, err = trans.SelectInt("SELECT COUNT(*) FROM EpicUserMap WHERE epic_id=?", epic_id);
        err == nil && count == 0 {
        epicToBeDeleted := Epic{Id: epic_id}
        if _, err = trans.Delete(&epicToBeDeleted); err == nil {
            if _, err = trans.Exec("DELETE FROM Story WHERE epic_id=?", epic_id); err == nil {
                trans.Commit()
            } else {
                utils.PrintErr(err, "removeUnownedEpic: Failed to delete stories for epic " +
                    strconv.FormatInt(epic_id, 10))
            }
        } else {
            utils.PrintErr(err, "removeUnownedEpic: Failed to delete epic " + strconv.FormatInt(epic_id, 10))
        }
    } else {
        utils.PrintErr(err, "removeUnownedEpic: Failed to select mappings for epic " + strconv.FormatInt(epic_id, 10))
    }
    trans.Rollback()
}
