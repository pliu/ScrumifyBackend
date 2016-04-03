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
	var check models.EpicModuleMap
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM EpicModuleMap WHERE moduleid=?", module_id)
	if err == nil && epicOwnedByUser(user_id, strconv.FormatInt(check.EpicID, 10)) {
		return true
	} else {
		return false
	}
}

func storyOwnedByUser(user_id string, story_id string) bool {
	var check models.ModuleStoryMap
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM ModuleStoryMap WHERE storyid=?", story_id)
	if err == nil && moduleOwnedByUser(user_id, strconv.FormatInt(check.ModuleID, 10)) {
		return true
	} else {
		return false
	}
}

func validUser() bool {
	return true
}

func validProject() bool {
	return true
}

func validModule() bool {
	return true
}

func validStory() bool {
	return true
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

// Need to check that
func putDependencies(module_id int64, dependencies []int64) {

}

func getDependencies(module_id int64) []int64 {
	var dependencies []int64
	return dependencies
}
