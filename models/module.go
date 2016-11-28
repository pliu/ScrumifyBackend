package models

import "gopkg.in/gorp.v2"

type Module struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	DueDate string `db:"duedate" json:"duedate"`
	Stage   int64  `db:"stage" json:"stage"`
	EpicId  int64  `db:"epidid" json:"epicid"`
}

func SetModuleProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.AddIndex("EpicIdIndex", "Hash", []string{"EpicId"})
	table.ColMap("Name").SetNotNull(true)
	table.ColMap("Stage").SetNotNull(true)
	table.ColMap("EpicId").SetNotNull(true)
}

func GetModule() {

}

func CreateModule() {

}

func UpdateModule() {

}

func DeleteModule() {

}
