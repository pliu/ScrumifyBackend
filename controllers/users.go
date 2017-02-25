package controllers

import (
	"ScrumifyBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"ScrumifyBackend/utils"
	"strconv"
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

	if !user.IsValid() {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	if user, err := models.CreateUpdateUser(user, false); err == nil {
		c.JSON(http.StatusCreated, user.Scrub())
	} else if err == utils.EmailExists {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var newUserInfo models.User
	c.Bind(&newUserInfo)

	if !newUserInfo.IsValid() {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	var int_id int64
	var err error
	if int_id, err = strconv.ParseInt(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestReturn)
		return
	}
	newUserInfo.Id = int_id
	if newUserInfo, err = models.CreateUpdateUser(newUserInfo, true); err == nil {
		c.JSON(http.StatusOK, newUserInfo.Scrub())
	} else if err == utils.EmailExists {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	if err := models.DeleteUser(id); err == nil {
		c.JSON(http.StatusOK, "User #" + id + " deleted")
	} else {
		c.JSON(http.StatusInternalServerError, utils.InternalErrorReturn)
	}
}
