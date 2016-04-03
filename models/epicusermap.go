package models

// Which users are part of which epics
type EpicUserMap struct {
	Id     int64 `db:"id" json:"id"`
	UserID int64 `db:"userid" json:"userid"`
	EpicID int64 `db:"epicid" json:"epicid"`
}

func GetEpicUserMap() {

}

func CreateEpicUserMap() {

}

func UpdateEpicUserMap() {

}

func DeleteEpicUserMap() {

}
