package models

import "gopkg.in/gorp.v2"

type Epic struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func SetEpicProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.ColMap("Name").SetNotNull(true)
}

func GetEpic() {

}

func CreateEpic() {

}

func UpdateEpic() {

}

func DeleteEpic() {

}

func (epic Epic)IsValid() bool {
	if epic.Name != "" {
		return true
	} else {
		return false
	}
}
