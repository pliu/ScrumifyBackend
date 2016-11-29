package models

import "gopkg.in/gorp.v2"

// Which users are part of which epics
type EpicUserMap struct {
	UserId int64 `db:"userid" json:"userid"`
	EpicId int64 `db:"epicid" json:"epicid"`
}

func SetEpicUserMapProperties(table *gorp.TableMap) {
	table.SetKeys(false, "UserId", "EpicId")
	table.AddIndex("UserIdIndex", "Hash", []string{"UserId"})
	table.AddIndex("EpicIdIndex", "Hash", []string{"EpicId"})
	table.SetUniqueTogether("UserId", "EpicId")
	table.ColMap("UserId").SetNotNull(true)
	table.ColMap("EpicId").SetNotNull(true)
}

func CreateEpicUserMap(mapping EpicUserMap) error {
	err := Dbmap.Insert(&mapping)
	return err
}

func DeleteEpicUserMap(mapping EpicUserMap) error {
	_, err := Dbmap.Delete(&mapping)
	return err
}

func EpicOwnedByUser(user_id string, epic_id string) (EpicUserMap, error) {
	var epicUserMap EpicUserMap
	err := Dbmap.SelectOne(&epicUserMap, "SELECT * FROM EpicUserMap WHERE userid=? AND epicid=?", user_id, epic_id)
	return epicUserMap, err
}
