package controllers

import (
	"ScrumifyBackend/models"
	"ScrumifyBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetStory(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	story_id := c.Params.ByName("storyid")

	if _, err := models.EpicOwnedByUser(user_id, epic_id); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if story, err := models.GetStory(epic_id, story_id); err != nil {
		if err == utils.StoryDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
	} else {
		c.JSON(http.StatusOK, story)
	}
}

func PostStory(c *gin.Context) {
	user_id := c.Params.ByName("id")
	var story models.Story
	c.Bind(&story)

	if !story.IsValid() {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	if _, err := models.EpicOwnedByUser(user_id, strconv.FormatInt(story.EpicId, 10)); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if assigneeError(story, c) {
		return
	}
	if story, err := models.CreateUpdateStory(story, false); err == nil {
		c.JSON(http.StatusCreated, story)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func UpdateStory(c *gin.Context) {
	user_id := c.Params.ByName("id")
	var story models.Story
	c.Bind(&story)

	if !story.IsValid() && story.Id != 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	if _, err := models.EpicOwnedByUser(user_id, strconv.FormatInt(story.EpicId, 10)); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if assigneeError(story, c) {
		return
	}
	if _, err := models.CreateUpdateStory(story, true); err == nil {
		c.JSON(http.StatusOK, story)
	} else if err == utils.StoryDoesntExist {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func DeleteStory(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	story_id := c.Params.ByName("storyid")

	if _, err := models.EpicOwnedByUser(user_id, epic_id); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if err := models.DeleteStory(epic_id, story_id); err == nil {
		c.JSON(http.StatusOK, "Story #" + story_id + " deleted")
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func assigneeError(story models.Story, c *gin.Context) bool {
	if story.AssignedTo > 0 {
		if _, err := models.EpicOwnedByUser(strconv.FormatInt(story.AssignedTo, 10), strconv.FormatInt(
			story.EpicId, 10)); err == utils.MappingDoesntExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User #" + strconv.FormatInt(story.AssignedTo, 10) +
					" not in epic #" + strconv.FormatInt(story.EpicId, 10)})
			return true
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
			return true
		}
	}
	return false
}
