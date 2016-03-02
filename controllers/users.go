package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	_, err := models.Dbmap.Select(&users, "SELECT * FROM User")

	if err == nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users in the table"})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &models.User{
			Id:       user_id,
			Username: user.Username,
			HashedPW: user.HashedPW,
			Email:    user.Email,
		}
		c.JSON(http.StatusOK, content)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func PostUser(c *gin.Context) {
	var user models.User
	c.Bind(&user)

	if user.Username != "" && user.HashedPW != "" && user.Email != "" {

		if insert, _ := models.Dbmap.Exec(`INSERT INTO User (username, hashedpw, email) VALUES (?, ?, ?)`, user.Username, user.HashedPW, user.Email); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				content := &models.User{
					Id:       user_id,
					Username: user.Username,
					HashedPW: user.HashedPW,
					Email:    user.Email,
				}
				c.JSON(201, content)
			} else {
				utils.CheckErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(422, gin.H{"error": "Field(s) is(are) empty"})
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		var json models.User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := models.User{
			Id:       user_id,
			Username: json.Username,
			HashedPW: json.HashedPW,
			Email:    json.Email,
		}

		if user.Username != "" && user.HashedPW != "" && user.Email != "" {
			_, err = models.Dbmap.Update(&user)

			if err == nil {
				c.JSON(200, user)
			} else {
				utils.CheckErr(err, "Updated failed")
			}

		} else {
			c.JSON(422, gin.H{"error": "Field(s) is(are) empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

// Still need to recursively delete all epics that contain only this user as a member
func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", id)

	if err == nil {
		_, err = models.Dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
		} else {
			utils.CheckErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}
