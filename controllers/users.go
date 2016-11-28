package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	users, err := models.GetUsers()

	if err == nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not access database"})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user, err := models.GetUser(id)

	if err == nil {
		c.JSON(http.StatusOK, models.ScrubUser(user))
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func PostUser(c *gin.Context) {
	var user models.User
	c.Bind(&user)

	if user.IsValid() {
		_, err := models.GetUserByEmail(user.Email)
		if err != nil {
			user, err = models.CreateUser(user)
			if err == nil {
				c.JSON(http.StatusCreated, user)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
	user, err := models.GetUser(id)

	if err == nil {
		var newUserInfo models.User
		c.Bind(&newUserInfo)
		existingUser, nerr := models.GetUserByEmail(newUserInfo.Email)
		if nerr != nil || existingUser.Id == user.Id {
			if newUserInfo.IsValid() {
				user := models.User{
					Id:       user.Id,
					Username: newUserInfo.Username,
					HashedPw: newUserInfo.HashedPw,
					Email:    newUserInfo.Email,
				}
				user, err = models.UpdateUser(user)
				if err == nil {
					c.JSON(http.StatusOK, user)
				} else {
					utils.PrintErr(err, "Update user failed")  // Should probably return something to the client
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already being used"})
		}

	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user, err := models.GetUser(id)

	if err == nil {
		err = models.DeleteUser(user)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"id #" + id: "User deleted"})
			go removeUserMappings(id)  // Need to come back to this
		} else {
			utils.PrintErr(err, "Delete user failed")
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}
