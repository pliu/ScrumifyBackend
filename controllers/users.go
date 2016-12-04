package controllers

import (
    "TodoBackend/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "TodoBackend/utils"
)

func GetUsers(c *gin.Context) {
    if users, err := models.GetUsers(); err == nil {
        c.JSON(http.StatusOK, users)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}

func GetUser(c *gin.Context) {
    id := c.Params.ByName("id")

    if user, err := models.GetUser(id); err == nil {
        c.JSON(http.StatusOK, user.Scrub())
    } else if err == utils.UserDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}

func PostUser(c *gin.Context) {
    var user models.User
    c.Bind(&user)

    if user.IsValid() {
        if user, err := models.CreateUpdateUser(user, false); err == nil {
            c.JSON(http.StatusCreated, user.Scrub())
        } else if err == utils.EmailExists {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else {
        c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
    }
}

func UpdateUser(c *gin.Context) {
    id := c.Params.ByName("id")
    var newUserInfo models.User
    c.Bind(&newUserInfo)

    if newUserInfo.IsValid() {

        // TODO: Authentication will make this step redundant
        if user, err := models.GetUser(id); err == nil {
            newUserInfo.Id = user.Id
            if newUserInfo, err = models.CreateUpdateUser(newUserInfo, true); err == nil {
                c.JSON(http.StatusOK, newUserInfo.Scrub())
            } else if err == utils.EmailExists {
                c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
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

func DeleteUser(c *gin.Context) {
    id := c.Params.ByName("id")

    // TODO: Authentication will make this step redundant
    if user, err := models.GetUser(id); err == nil {
        if err = models.DeleteUser(user); err == nil {
            c.JSON(http.StatusOK, "User #" + id + " deleted")
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else if err == utils.UserDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}
