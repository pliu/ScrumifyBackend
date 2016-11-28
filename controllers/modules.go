package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetModules(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	if epicOwnedByUser(id, epic_id) {
		var modules []models.RestModule
		_, err := models.Dbmap.Select(&modules, "SELECT * FROM Module WHERE epicid=?", epic_id)

		for _, module := range modules {
			module.Dependencies = getDependencies(module.Id)
		}

		if err == nil {
			c.JSON(http.StatusOK, modules)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not access database"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}

func PostModule(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	if epicOwnedByUser(id, epic_id) {
		var module models.RestModule
		c.Bind(&module)

		if validModule(module, epic_id) {

			if insert, _ := models.Dbmap.Exec(`INSERT INTO Module (name, duedate, stage, epicid) VALUES (?, ?, ?, ?)`, module.Name, module.DueDate, module.Stage, epic_id); insert != nil {
				module_id, err := insert.LastInsertId()
				if err == nil {
					putDependencies(module_id, module.Dependencies)
					module.Id = module_id
					c.JSON(http.StatusCreated, module)
				} else {
					utils.PrintErr(err, "Insert module failed")
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

	if moduleOwnedByUser(id, module_id) {
		var module models.Module
		err := models.Dbmap.SelectOne(&module, "SELECT * FROM Module WHERE id=?", module_id)

		if err == nil {

			_, err := models.Dbmap.Delete(&module)

			if err == nil {
				c.JSON(http.StatusOK, gin.H{"id #" + module_id: "Deleted module"})
				go removeModuleStories(module.Id)
				go removeModuleDependencies(module.Id)
			} else {
				utils.PrintErr(err, "Delete module failed")
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Module not owned by you"})
	}
}
