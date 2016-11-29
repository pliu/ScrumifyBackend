// Library of functions to asynchronously clean up after deletions

package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
)

// Async
func removeUserMappings(user_id string) {
	var mappings []models.EpicUserMap
	models.Dbmap.Select(&mappings, "SELECT * FROM EpicUserMap WHERE userid=?", user_id)
	for _, mapping := range mappings {
		_, err := models.Dbmap.Delete(&mapping)
		if err != nil {
			utils.PrintErr(err, "Delete user mapping failed")
		}
		removeUnownedEpic(mapping.EpicId)
	}
}

// Async
func removeUnownedEpic(epic_id int64) {
	var mappings []models.EpicUserMap
	_, err := models.Dbmap.Select(&mappings, "SELECT * FROM EpicUserMap WHERE epicid=?", epic_id)
	if len(mappings) == 0 {
		epic := models.Epic{
			Id: epic_id,
		}
		_, err = models.Dbmap.Delete(&epic)
		if err != nil {
			utils.PrintErr(err, "Delete unowned epic failed")
		} else {
			removeEpicStories(epic_id)
		}
	}
}

// Async
func removeEpicStories(epic_id int64) {
	var stories []models.Story
	_, err := models.Dbmap.Select(&stories, "SELECT * FROM Module WHERE epicid=?", epic_id)
	for _, story := range stories {
		_, err = models.Dbmap.Delete(&story)
		if err != nil {
			utils.PrintErr(err, "Delete module failed")
		} else {

		}
	}
}
