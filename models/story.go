package models

import (
    "gopkg.in/gorp.v2"
    "TodoBackend/utils"
)

type Story struct {
    Id          int64  `db:"id" json:"id"`
    Name        string `db:"name" json:"name"`
    Description string `db:"description" json:"description"`
    DueDate     string `db:"due_date" json:"due_date"`
    Points      int64  `db:"points" json:"points"`
    Stage       int64  `db:"stage" json:"stage"`
    EpicId      int64  `db:"epic_id" json:"epic_id"`
}

func SetStoryProperties(table *gorp.TableMap) {
    table.SetKeys(true, "Id")

    // InnoDB does not have Hash indices
    table.AddIndex("StoryEpicIdIndex", "Hash", []string{"EpicId"})
    table.ColMap("Name").SetNotNull(true)
    table.ColMap("Stage").SetNotNull(true)
    table.ColMap("EpicId").SetNotNull(true)
}

func GetStories(epic_id string) ([]Story, error) {
    var stories []Story
    _, err := Dbmap.Select(&stories, "SELECT * FROM Story WHERE epic_id=?", epic_id)
    utils.PrintErr(err, "GetStories: Failed to select stories for epic " + epic_id)
    return stories, err
}

func CreateStory() {

}

func UpdateStory() {

}

func DeleteStory() {

}

func (story Story)IsValid() bool {
    if story.Name != "" && (story.Stage == 0 || story.Stage == 1 || story.Stage == 2) {
        return true
    } else {
        return false
    }
}
