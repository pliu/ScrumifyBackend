package controllers

import (
	"TodoBackend/models"
	"strconv"
)

func storyOwnedByUser(user_id string, story_id string) bool {
	var check models.Story
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM Story WHERE id=?", story_id)
	if err == nil && models.EpicOwnedByUser(user_id, strconv.FormatInt(check.EpicId, 10)) {
		return true
	} else {
		return false
	}
}

func validStory(story models.Story) bool {
	if story.Name != "" && (story.Stage == 0 || story.Stage == 1 || story.Stage == 2) {
		return true
	} else {
		return false
	}
}

func putDependencies(module_id int64, dependencies []int64) {

}
