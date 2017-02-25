package controllers

import (
	"ScrumifyBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"ScrumifyBackend/utils"
)

func GetEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if _, err := models.EpicOwnedByUser(user_id, epic_id); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if epic, err := models.GetEpic(epic_id); err == nil {
		for i := range epic.Members {
			epic.Members[i].HashedPw = ""
		}
		c.JSON(http.StatusOK, epic)
	} else if err == utils.EpicDoesntExist {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func PostEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	var epic models.Epic
	c.Bind(&epic)

	if !epic.IsValid() {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	if epic, err := models.CreateEpic(user_id, epic); err == nil {
		c.JSON(http.StatusCreated, epic)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func UpdateEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	var newEpicInfo models.Epic
	c.Bind(&newEpicInfo)

	if !newEpicInfo.IsValid() {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	var err error
	if _, err = models.EpicOwnedByUser(user_id, strconv.FormatInt(newEpicInfo.Id, 10)); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if newEpicInfo, err = models.UpdateEpic(newEpicInfo); err == nil {
		c.JSON(http.StatusOK, newEpicInfo)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func DeleteEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	var mapping models.EpicUserMap
	var err error
	if mapping, err = models.EpicOwnedByUser(user_id, epic_id); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if err = models.DeleteEpic(mapping); err == nil {
		c.JSON(http.StatusOK, "Deleted epic #" + epic_id + " from user #" + user_id + "'s list")
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func AddUserToEpic(c *gin.Context) {
	user_id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")
	var email models.RestEmail
	c.Bind(&email)

	var mapping models.EpicUserMap
	var err error
	if mapping, err = models.EpicOwnedByUser(user_id, epic_id); err != nil {
		if err == utils.MappingDoesntExist {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
		} else {
			c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
		}
		return
	}
	if user, err := models.AddEpicUserMap(email.Email, epic_id); err == nil {
		c.JSON(http.StatusOK, "User #" + strconv.FormatInt(user.Id, 10) + " associated with epic #" + epic_id)
	} else if err == utils.CantParseEpicId {
		c.JSON(http.StatusBadRequest, gin.H{"error": epic_id + " is not a valid epic ID"})
	} else if err == utils.EmailDoesntExist {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if err == utils.MappingExists {
		c.JSON(http.StatusOK, "User #" + strconv.FormatInt(mapping.UserId, 10) + " already associated with epic #" +
				epic_id)
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}
