package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

		c.JSON(http.StatusOK, scrubUser(user))
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func PostUser(c *gin.Context) {
	var user models.User
	c.Bind(&user)

	if user.Username != "" && user.HashedPW != "" && user.Email != "" {

		user.Email = strings.ToUpper(user.Email)
		if insert, _ := models.Dbmap.Exec(`INSERT INTO User (username, hashedpw, email) VALUES (?, ?, ?)`, user.Username, user.HashedPW, user.Email); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				user.Id = user_id
				user.Email = user.Email
				c.JSON(201, scrubUser(user))
			} else {
				utils.CheckErr(err, "Insert user failed")
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

		user := models.User{
			Id:       user.Id,
			Username: json.Username,
			HashedPW: json.HashedPW,
			Email:    strings.ToUpper(json.Email),
		}

		if user.Username != "" && user.HashedPW != "" && user.Email != "" {
			_, err = models.Dbmap.Update(&user)

			if err == nil {
				c.JSON(200, scrubUser(user))
			} else {
				utils.CheckErr(err, "Update user failed")
			}

		} else {
			c.JSON(422, gin.H{"error": "Field(s) is(are) empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", id)

	if err == nil {
		_, err = models.Dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
			removeUserMappings(id)
		} else {
			utils.CheckErr(err, "Delete user failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func GetUserByEmail(email string) (models.User, error) {
	email = strings.ToUpper(email)
	var user models.User
	err := models.Dbmap.SelectOne(&user, "SELECT id FROM User WHERE email=?", email)
	return user, err
}

func removeUserMappings(user_id string) {
	var mappings []models.EpicUserMap
	models.Dbmap.Select(&mappings, "SELECT * FROM EpicUserMap WHERE userid=?", user_id)
	for _, mapping := range mappings {
		_, err := models.Dbmap.Delete(&mapping)
		if err != nil {
			utils.CheckErr(err, "Delete user mapping failed")
		}
		RemoveUnownedEpic(mapping.EpicID)
	}
}

func scrubUser(user models.User) models.User {
	scrubbed_user := models.User{
		Id:       user.Id,
		Username: user.Username,
		HashedPW: "",
		Email:    user.Email,
	}
	return scrubbed_user
}
