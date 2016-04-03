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
			utils.CheckErr(err, "Delete user mapping failed")
		}
		removeUnownedEpic(mapping.EpicID)
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
			utils.CheckErr(err, "Delete unowned epic failed")
		} else {
			removeEpicModules(epic_id)
		}
	}
}

// Async
func removeEpicModules(epic_id int64) {
	var mappings []models.EpicModuleMap
	_, err := models.Dbmap.Select(&mappings, "SELECT * FROM EpicModuleMap WHERE epicid=?", epic_id)
	for _, mapping := range mappings {
		module := models.Module{
			Id: mapping.ModuleID,
		}
		_, err = models.Dbmap.Delete(&module)
		if err != nil {
			utils.CheckErr(err, "Delete module failed")
		} else {
			models.Dbmap.Delete(&mapping)
			removeModuleStories(module.Id)
			removeModuleDependencies(module.Id)
		}
	}
}

// Async
func removeModuleStories(module_id int64) {
	var mappings []models.ModuleStoryMap
	_, err := models.Dbmap.Select(&mappings, "SELECT * FROM ModuleStoryMap WHERE moduleid=?", module_id)
	for _, mapping := range mappings {
		story := models.Story{
			Id: mapping.StoryID,
		}
		_, err = models.Dbmap.Delete(&story)
		if err != nil {
			utils.CheckErr(err, "Delete story failed")
		} else {
			models.Dbmap.Delete(&mapping)
		}
	}
}

// Async
func removeModuleDependencies(module_id int64) {

}
