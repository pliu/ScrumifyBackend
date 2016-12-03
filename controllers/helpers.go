package controllers

import (
    "TodoBackend/models"
    "strconv"
)

func storyOwnedByUser(user_id string, story_id string) bool {
    var check models.Story
    if err := models.Dbmap.SelectOne(&check, "SELECT * FROM Story WHERE id=?", story_id); err == nil {
        if _, err = models.EpicOwnedByUser(user_id, strconv.FormatInt(check.EpicId, 10)); err == nil {
            return true
        }
    }
    return false
}
