package controllers

import (
    "ScrumifyBackend/models"
    "ScrumifyBackend/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func GetStories(c *gin.Context) {
    user_id := c.Params.ByName("id")
    epic_id := c.Params.ByName("epicid")

    if _, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
        if modules, err := models.GetStories(epic_id); err == nil {
            c.JSON(http.StatusOK, modules)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    }
}

func PostStory(c *gin.Context) {
    user_id := c.Params.ByName("id")
    var story models.Story
    c.Bind(&story)

    if story.IsValid() {
        if _, err := models.EpicOwnedByUser(user_id, strconv.FormatInt(story.EpicId, 10)); err == nil {
            if story, err = models.CreateUpdateStory(story, false); err == nil {
                c.JSON(http.StatusCreated, story)
            } else {
                c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
            }
        } else if err == utils.MappingDoesntExist {
            c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else {
        c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
    }
}

func UpdateStory(c *gin.Context) {
	user_id := c.Params.ByName("id")
	var newStoryInfo models.Story
    c.Bind(&newStoryInfo)

    if newStoryInfo.IsValid() {
        if story, err := models.GetStory(strconv.FormatInt(newStoryInfo.Id, 10)); err == nil {
            if _, err = models.EpicOwnedByUser(user_id, strconv.FormatInt(story.EpicId, 10)); err == nil {
                newStoryInfo.EpicId = story.EpicId
                if newStoryInfo, err = models.CreateUpdateStory(newStoryInfo, true); err == nil {
                    c.JSON(http.StatusOK, newStoryInfo)
                } else {
                    c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
                }
            } else if err == utils.MappingDoesntExist {
                c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
            } else {
                c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
            }
        } else if err == utils.StoryDoesntExist {
            c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else {
        c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
    }
}

func DeleteStory(c *gin.Context) {
    user_id := c.Params.ByName("id")
    story_id := c.Params.ByName("storyid")

    if story, err := models.GetStory(story_id); err == nil {
        if _, err = models.EpicOwnedByUser(user_id, strconv.FormatInt(story.EpicId, 10)); err == nil {
            if err = models.DeleteStory(story); err == nil {
                c.JSON(http.StatusOK, "Story #" + story_id + " deleted")
            } else {
                c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
            }
        } else if err == utils.MappingDoesntExist {
            c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else if err == utils.StoryDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}
