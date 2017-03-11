package models

import (
	"gopkg.in/gorp.v2"
	"ScrumifyBackend/utils"
	"strconv"
	"time"
	"database/sql"
)

type Epic struct {
	Id        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Members   []User    `db:"-" json:"members"`
	Stories   []Story   `db:"-" json:"stories"`
	CreatedAt time.Time `db:"created_at" json:"-"`
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

func GetEpic(epic_id string) (Epic, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "GetEpic: Failed to begin transaction")
		return Epic{}, err
	}

	var epic Epic
	if err := trans.SelectOne(&epic, "SELECT * FROM Epic WHERE id=?", epic_id); err != nil {
		if err == sql.ErrNoRows {
			trans.Rollback()
			return Epic{}, utils.EpicDoesntExist
		} else {
			trans.Rollback()
			utils.PrintErr(err, "GetEpic: Failed to select epic " + epic_id)
			return Epic{}, err
		}
	}
	if _, err = trans.Select(&epic.Stories, "SELECT * FROM Story WHERE epic_id=?", epic_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetEpic: Failed to select stories for epic " + epic_id)
		return Epic{}, err
	}
	if _, err = trans.Select(&epic.Members, "SELECT * FROM User WHERE id IN (SELECT user_id FROM EpicUserMap WHERE " +
			"epic_id=?)", epic_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetEpic: Failed to select members for epic " + epic_id)
		return Epic{}, err
	}
	return epic, trans.Commit()
}

func CreateEpic(user_id string, epic Epic) (Epic, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "CreateEpic: Failed to begin transaction")
		return Epic{}, err
	}

	if err = trans.Insert(&epic); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "CreateEpic: Failed to insert epic")
		return Epic{}, err
	}
	if _, err = trans.Exec("INSERT INTO EpicUserMap (user_id, epic_id) VALUES (?, ?)", user_id, epic.Id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "CreateEpic: Failed to insert mapping for user_id " + user_id)
		return Epic{}, err
	}
	var check Epic
	if err = trans.SelectOne(&check, "SELECT * FROM Epic WHERE id=?", epic.Id); err == nil {
		return check, trans.Commit()
	} else {
		return epic, trans.Commit()
	}
}

func UpdateEpic(epic Epic) (Epic, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "UpdateEpic: Failed to begin transaction")
		return Epic{}, err
	}

	if _, err = trans.Update(&epic); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "UpdateEpic: Failed to update epic " + strconv.FormatInt(epic.Id, 10))
		return Epic{}, err
	}
	return epic, trans.Commit()
}

func DeleteEpic(mapping EpicUserMap) error {
	_, err := Dbmap.Delete(&mapping)
	utils.PrintErr(err, "DeleteEpic: Failed to delete mapping for user_id " + strconv.FormatInt(mapping.UserId, 10) +
			" and epic_id " + strconv.FormatInt(mapping.EpicId, 10))
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

	if _, err = trans.Exec("DELETE FROM Epic WHERE id NOT IN (SELECT DISTINCT epic_id FROM EpicUserMap WHERE " +
			"epic_id=?)", epic_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "removeUnownedEpic: Failed to delete epic " + strconv.FormatInt(epic_id, 10))
	} else {
		trans.Commit()
	}
}
