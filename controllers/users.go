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
        c.JSON(http.StatusOK, scrubUser(user))
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
            c.JSON(http.StatusCreated, scrubUser(user))
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
    var user models.User
    c.Bind(&user)

    if user.IsValid() {
        if dummy, err := models.GetUser(id); err == nil {
            // TODO: Authentication will make this step redundant
            user.Id = dummy.Id
            if user, err = models.CreateUpdateUser(user, true); err == nil {
                c.JSON(http.StatusOK, scrubUser(user))
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

    if user, err := models.GetUser(id); err == nil {
        // TODO: Authentication will make this step redundant
        if err = models.DeleteUser(user); err == nil {
            c.JSON(http.StatusOK, gin.H{"id #" + id: "User deleted"})
        } else {
            c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
        }
    } else if err == utils.UserDoesntExist {
        c.JSON(http.StatusUnauthorized, utils.UnauthorizedReturn)
    } else {
        c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
    }
}

func scrubUser(user models.User) models.User {
    return models.User{
        Id:       user.Id,
        Username: user.Username,
        HashedPw: "",
        Email:    user.Email,
    }
}
