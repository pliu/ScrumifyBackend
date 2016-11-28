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

func GetEpicUserMap() {

}

func CreateEpicUserMap() {

}

func UpdateEpicUserMap() {

}

func DeleteEpicUserMap() {

}
