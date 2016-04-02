package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetModules(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	if CheckEpicOwnedByUser(id, epic_id) {
		var modules []models.ModuleIn
		_, err := models.Dbmap.Select(&modules, "SELECT * FROM Module WHERE id IN (SELECT moduleid FROM EpicModuleMap WHERE epicid=?)", epic_id)

		for _, module := range modules {
			module.Dependencies = getDependencies(module.Id)
		}

		if err == nil {
			c.JSON(http.StatusOK, modules)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Can't find associated modules"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}

func PostModule(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	if CheckEpicOwnedByUser(id, epic_id) {
		var module models.ModuleIn
		c.Bind(&module)

		if module.Name != "" && (module.Stage == 0 || module.Stage == 1 || module.Stage == 2) {

			if insert, _ := models.Dbmap.Exec(`INSERT INTO Module (name, duedate, stage) VALUES (?, ?, ?)`, module.Name, module.DueDate, module.Stage); insert != nil {
				module_id, err := insert.LastInsertId()
				if err == nil {
					models.Dbmap.Exec(`INSERT INTO EpicModuleMap (moduleid, epicid) VALUES (?, ?)`, module_id, epic_id)
					putDependencies(module_id, module.Dependencies)
					module.Id = module_id
					c.JSON(http.StatusCreated, module)
				} else {
					utils.CheckErr(err, "Insert module failed")
				}
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}

/*
func UpdateModule(c *gin.Context) {
	id := c.Params.ByName("id")
	module_id := c.Params.ByName("moduleid")

	if if CheckModuleOwnedByUser(id, module_id) {
		var module models.ModuleIn
		err := models.Dbmap.SelectOne(&module, "SELECT * FROM Module WHERE id=?", module_id)
		module.Dependencies = getDependencies(module.Id)

		if err == nil {
			var json models.ModuleIn
			c.Bind(&json)

			module := models.Epic{
				Id:      module.Id,
				Name:    json.Name,
				DueDate: json.DueDate,
				Stage:   json.Stage,
			}

			if epic.Name != "" {
				_, err = models.Dbmap.Update(&epic)

				if err == nil {
					c.JSON(200, epic)
				} else {
					utils.CheckErr(err, "Updated epic failed")
				}

			} else {
				c.JSON(422, gin.H{"error": "Field(s) is(are) empty"})
			}

		} else {
			c.JSON(404, gin.H{"error": "Module not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic/module not owned by you"})
	}
}
*/

func DeleteModule(c *gin.Context) {
	id := c.Params.ByName("id")
	module_id := c.Params.ByName("moduleid")

	if CheckModuleOwnedByUser(id, module_id) {
		var module models.Module
		err := models.Dbmap.SelectOne(&module, "SELECT * FROM Module WHERE id=?", module_id)

		if err == nil {

			_, err := models.Dbmap.Delete(&module)

			if err == nil {
				c.JSON(http.StatusOK, gin.H{"id #" + module_id: "Deleted module"})
				var mapping models.EpicModuleMap
				models.Dbmap.SelectOne(&mapping, "SELECT * FROM EpicModuleMap WHERE moduleid=?", module_id)
				models.Dbmap.Delete(&mapping)
				RemoveModuleStories(module.Id)
				removeModuleDependencies(module.Id)
			} else {
				utils.CheckErr(err, "Delete module failed")
			}

		} else {
			c.JSON(404, gin.H{"error": "Module not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Module not owned by you"})
	}
}

func RemoveEpicModules(epic_id int64) {
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
			RemoveModuleStories(module.Id)
			removeModuleDependencies(module.Id)
		}
	}
}

func CheckModuleOwnedByUser(user_id string, module_id string) bool {
	var check models.EpicModuleMap
	err := models.Dbmap.SelectOne(&check, "SELECT * FROM EpicModuleMap WHERE moduleid=?", module_id)
	if err == nil && CheckEpicOwnedByUser(user_id, strconv.FormatInt(check.EpicID, 10)) {
		return true
	} else {
		return false
	}
}

func putDependencies(module_id int64, dependencies []int64) {

}

func getDependencies(module_id int64) []int64 {
	var dependencies []int64
	return dependencies
}

func removeModuleDependencies(module_id int64) {

}
