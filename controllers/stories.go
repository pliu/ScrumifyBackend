package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStories(c *gin.Context) {
	id := c.Params.ByName("id")
	module_id := c.Params.ByName("moduleid")
	if moduleOwnedByUser(id, module_id) {
		var stories []models.Story
		_, err := models.Dbmap.Select(&stories, "SELECT * FROM Story WHERE owner=?", module_id)

		if err == nil {
			c.JSON(http.StatusOK, stories)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not access database"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Module not owned by you"})
	}
}

func PostStory(c *gin.Context) {
	id := c.Params.ByName("id")
	module_id := c.Params.ByName("moduleid")
	if moduleOwnedByUser(id, module_id) {
		var story models.Story
		c.Bind(&story)

		if validStory(story) {

			if insert, _ := models.Dbmap.Exec(`INSERT INTO Story (name, stage, description, points, owner) VALUES (?, ?, ?, ?, ?)`, story.Name, story.Stage, story.Description, story.Points, module_id); insert != nil {
				story_id, err := insert.LastInsertId()
				if err == nil {
					story.Id = story_id
					c.JSON(http.StatusCreated, story)
				} else {
					utils.CheckErr(err, "Insert story failed")
				}
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Module not owned by you"})
	}
}

/*
func UpdateStory(c *gin.Context) {
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

func DeleteStory(c *gin.Context) {
	id := c.Params.ByName("id")
	story_id := c.Params.ByName("storyid")

	if storyOwnedByUser(id, story_id) {
		var story models.Story
		err := models.Dbmap.SelectOne(&story, "SELECT * FROM Story WHERE id=?", story_id)

		if err == nil {

			_, err := models.Dbmap.Delete(&story)

			if err == nil {
				c.JSON(http.StatusOK, gin.H{"id #" + story_id: "Deleted story"})
			} else {
				utils.CheckErr(err, "Delete story failed")
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Story not owned by you"})
	}
}
