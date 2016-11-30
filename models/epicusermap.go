package models

import "gopkg.in/gorp.v2"

// Which users are part of which epics
type EpicUserMap struct {
	UserId int64 `db:"userid" json:"userid"`
	EpicId int64 `db:"epicid" json:"epicid"`
}

func SetEpicUserMapProperties(table *gorp.TableMap) {
	table.SetKeys(false, "UserId", "EpicId")

	// InnoDB does not have Hash indices
	table.AddIndex("MapUserIdIndex", "Btree", []string{"UserId"})
	table.AddIndex("MapEpicIdIndex", "Btree", []string{"EpicId"})
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
