package models

import "gopkg.in/gorp.v2"

type Story struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	DueDate 	string  `db:"duedate" json:"duedate"`
	Points      int64  `db:"points" json:"points"`
	Stage       int64  `db:"stage" json:"stage"`
	EpicId    	int64  `db:"epicid" json:"epicid"`
}

func SetStoryProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")
	table.AddIndex("ModuleIdIndex", "Hash", []string{"ModuleId"})
	table.ColMap("Name").SetNotNull(true)
	table.ColMap("Stage").SetNotNull(true)
	table.ColMap("EpicId").SetNotNull(true)
}

func GetStories(epic_id string) ([]Story, error) {
	var stories []Story
	_, err := Dbmap.Select(&stories, "SELECT * FROM Module WHERE epicid=?", epic_id)
	return stories, err
}

func CreateStory() {

}

func UpdateStory() {

}

func DeleteStory() {

}
