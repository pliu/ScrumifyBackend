package controllers

import (
	"TodoBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetEpics(c *gin.Context) {
	user_id := c.Params.ByName("id")

	if epics, err := models.GetEpics(user_id); err == nil {
		c.JSON(http.StatusOK, epics)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func PostEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")

	if _, err := models.GetUser(user_id); err == nil {
		var epic models.Epic
		c.Bind(&epic)

		if epic.IsValid() {
			if createdEpic, err := models.CreateEpic(user_id, epic); err == nil {
				c.JSON(http.StatusCreated, createdEpic)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func UpdateEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if _, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
		if existingEpic, err := models.GetEpic(epic_id); err == nil {
			var newEpicInfo models.Epic
			c.Bind(&newEpicInfo)

			if newEpicInfo.IsValid() {
				epic := models.Epic{
					Id:   existingEpic.Id,
					Name: newEpicInfo.Name,
				}

				if err = models.UpdateEpic(epic); err == nil {
					c.JSON(http.StatusOK, epic)
				} else {
					c.JSON(http.StatusInternalServerError, err.Error())
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
}

func DeleteEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if mapping, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
		if err = models.DeleteEpicUserMap(mapping); err == nil {
			c.JSON(http.StatusOK, gin.H{"id #" + epic_id: "Deleted from " + user_id + "'s list"})
			go removeUnownedEpic(mapping.EpicId)  // Need to come back to this
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
}

func AddUserToEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if _, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
		var email models.RestEmail
		c.Bind(&email)

		if user, nerr := models.GetUserByEmail(email.Email); nerr == nil {
			if _, err = models.EpicOwnedByUser(strconv.FormatInt(user.Id, 10), epic_id); err != nil {
				int_epic_id, _ := strconv.ParseInt(epic_id, 10, 64)
				int_user_id, _ := strconv.ParseInt(user_id, 10, 64)
				mapping := models.EpicUserMap{
					EpicId:	int_epic_id,
					UserId: int_user_id,
				}
				if err = models.CreateEpicUserMap(mapping); err == nil {
					c.JSON(http.StatusOK, strconv.FormatInt(user.Id, 10) + " associated with " + epic_id)
				} else {
					c.JSON(http.StatusInternalServerError, err.Error())
				}
			} else {
				c.JSON(http.StatusOK, gin.H{"error": "User already a member of the epic"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
}
