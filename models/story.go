package models

import "gopkg.in/gorp.v2"

type Story struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Points      int64  `db:"points" json:"points"`
	Stage       int64  `db:"stage" json:"stage"`
	ModuleId    int64  `db:"moduleid" json:"moduleid"`
}

func SetStoryProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.AddIndex("ModuleIdIndex", "Hash", []string{"ModuleId"})
	table.ColMap("Name").SetNotNull(true)
	table.ColMap("Stage").SetNotNull(true)
	table.ColMap("ModuleId").SetNotNull(true)
}

func GetStory() {

}

func CreateStory() {

}

func UpdateStory() {

}

func DeleteStory() {

}
