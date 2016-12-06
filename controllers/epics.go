package controllers

import (
    "ScrumifyBackend/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "ScrumifyBackend/utils"
)

func GetEpics(c *gin.Context) {
    user_id := c.Params.ByName("id")

    if _, err := models.GetUser(user_id); err == nil {
        if epics, err := models.GetEpics(user_id); err == nil {
            c.JSON(http.StatusOK, epics)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else if err == utils.UserDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}

func PostEpic(c *gin.Context) {
    user_id := c.Params.ByName("id")
    var epic models.Epic
    c.Bind(&epic)

    if epic.IsValid() {

        // TODO: Authentication will make this step redundant
        if _, err := models.GetUser(user_id); err == nil {
            if epic, err = models.CreateEpic(user_id, epic); err == nil {
                c.JSON(http.StatusCreated, epic)
            } else {
                c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
            }
        } else if err == utils.UserDoesntExist {
            c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else {
        c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
    }
}

func UpdateEpic(c *gin.Context) {
    user_id := c.Params.ByName("id")
    var newEpicInfo models.Epic
    c.Bind(&newEpicInfo)

    if newEpicInfo.IsValid() {
        if mapping, err := models.EpicOwnedByUser(user_id, strconv.FormatInt(newEpicInfo.Id, 10)); err == nil {
            newEpicInfo.Id = mapping.EpicId
            if newEpicInfo, err = models.UpdateEpic(newEpicInfo); err == nil {
                c.JSON(http.StatusOK, newEpicInfo)
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

func DeleteEpic(c *gin.Context) {
    user_id := c.Params.ByName("id")
    epic_id := c.Params.ByName("epicid")

    if mapping, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
        if err = models.DeleteEpic(mapping); err == nil {
            c.JSON(http.StatusOK, "Deleted epic #" + epic_id + " from user #" + user_id + "'s list")
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else if err == utils.MappingDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}

func AddUserToEpic(c *gin.Context) {
    user_id := c.Params.ByName("id")
    epic_id := c.Params.ByName("epicid")
    var email models.RestEmail
    c.Bind(&email)

    if mapping, err := models.EpicOwnedByUser(user_id, epic_id); err == nil {
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
    } else if err == utils.MappingDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}
