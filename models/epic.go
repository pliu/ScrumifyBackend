package models

import (
	"gopkg.in/gorp.v2"
)

type Epic struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func SetEpicProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.ColMap("Name").SetNotNull(true)
}

func GetEpics(user_id string) ([]Epic, error) {
	var epics []Epic
	_, err := Dbmap.Select(&epics, "SELECT * FROM Epic WHERE id IN (SELECT epicid FROM EpicUserMap WHERE userid=?)",
		user_id)
	return epics, err
}

func GetEpic(epic_id string) (Epic, error) {
	var epic Epic
	err := Dbmap.SelectOne(&epic, "SELECT * FROM Epic WHERE id=?", epic_id)
	return epic, err
}

func CreateEpic(user_id string, epic Epic) (Epic, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		return epic, err
	}

	if err := trans.Insert(&epic); err == nil {
		if _, err = trans.Exec(`INSERT INTO EpicUserMap (userid, epicid) VALUES (?, ?)`, user_id, epic.Id); err == nil {
			return epic, trans.Commit()
		} else {
			trans.Rollback();
			return epic, err
		}
	} else {
		trans.Rollback()
		return epic, err
	}
}

func UpdateEpic(epic Epic) error {
	_, err := Dbmap.Update(&epic)
	return err
}

func DeleteEpic(epic Epic) error {
	_, err := Dbmap.Delete(&epic)
	return err
}

func (epic Epic)IsValid() bool {
	if epic.Name != "" {
		return true
	} else {
		return false
	}
}
