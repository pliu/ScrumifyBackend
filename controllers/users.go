package controllers

import (
	"TodoBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	if users, err := models.GetUsers(); err == nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")

	if user, err := models.GetUser(id); err == nil {
		c.JSON(http.StatusOK, models.ScrubUser(user))
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func PostUser(c *gin.Context) {
	var user models.User
	c.Bind(&user)

	if user.IsValid() {
		if _, err := models.GetUserByEmail(user.Email); err != nil {
			if user, err = models.CreateUser(user); err == nil {
				c.JSON(http.StatusCreated, user)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already being used"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")

	if user, err := models.GetUser(id); err == nil {
		var newUserInfo models.User
		c.Bind(&newUserInfo)

		if existingUser, nerr := models.GetUserByEmail(newUserInfo.Email); nerr != nil || existingUser.Id == user.Id {
			if newUserInfo.IsValid() {
				user := models.User{
					Id:       user.Id,
					Username: newUserInfo.Username,
					HashedPw: newUserInfo.HashedPw,
					Email:    newUserInfo.Email,
				}
				if user, err = models.UpdateUser(user); err == nil {
					c.JSON(http.StatusOK, user)
				} else {
					c.JSON(http.StatusInternalServerError, err.Error())
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already being used"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	if user, err := models.GetUser(id); err == nil {
		if err = models.DeleteUser(user); err == nil {
			c.JSON(http.StatusOK, gin.H{"id #" + id: "User deleted"})
			go removeUserMappings(id)  // Need to come back to this
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}
