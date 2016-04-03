package controllers

import (
	"TodoBackend/models"
	"strconv"
	"strings"
)

func userExists(user_id string) bool {
	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", user_id)
	if err == nil {
		return true
	} else {
		return false
	}
}

func epicOwnedByUser(user_id string, epic_id string) bool {
	var check models.EpicUserMap
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM EpicUserMap WHERE userid=? AND epicid=?", user_id, epic_id)
	if err == nil && userExists(user_id) {
		return true
	} else {
		return false
	}
}

func moduleOwnedByUser(user_id string, module_id string) bool {
	var check models.Module
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM Module WHERE id=?", module_id)
	if err == nil && epicOwnedByUser(user_id, strconv.FormatInt(check.Owner, 10)) {
		return true
	} else {
		return false
	}
}

func storyOwnedByUser(user_id string, story_id string) bool {
	var check models.Story
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM Story WHERE id=?", story_id)
	if err == nil && moduleOwnedByUser(user_id, strconv.FormatInt(check.Owner, 10)) {
		return true
	} else {
		return false
	}
}

func validUser(user models.User) bool {
	if user.Username != "" && user.HashedPW != "" && user.Email != "" {
		return true
	} else {
		return false
	}
}

func validEpic(epic models.Epic) bool {
	if epic.Name != "" {
		return true
	} else {
		return false
	}
}

func validModule(module models.RestModule, epic_id string) bool {
	if module.Name != "" && (module.Stage == 0 || module.Stage == 1 || module.Stage == 2) && validDependencies(module.Dependencies, epic_id) {
		return true
	} else {
		return false
	}
}

func validDependencies(dependencies []int64, epic_id string) bool {
	return true
}

func moduleInEpic(module_id int64, epic_id string) bool {
	return true
}

func validStory(story models.Story) bool {
	if story.Name != "" && (story.Stage == 0 || story.Stage == 1 || story.Stage == 2) {
		return true
	} else {
		return false
	}
}

func getUserByEmail(email string) (models.User, error) {
	email = strings.ToUpper(email)
	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT id FROM User WHERE email=?", email)
	return user, err
}

func scrubUser(user models.User) models.User {
	scrubbed_user := models.User{
		Id:       user.Id,
		Username: user.Username,
		HashedPW: "",
		Email:    user.Email,
	}
	return scrubbed_user
}

func getDependencies(module_id int64) []int64 {
	var dependencies []int64
	return dependencies
}

func putDependencies(module_id int64, dependencies []int64) {

}
