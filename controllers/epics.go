package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	//"strconv"
)

func GetEpics(c *gin.Context) {
	id := c.Params.ByName("id")
	epics, err := GetEpicsByUser(id)

	if err == nil {
		c.JSON(http.StatusOK, epics)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can't find associated epics"})
	}
}

func GetEpicsByUser(user_id string) ([]models.Epic, error) {
	var epics []models.Epic
	_, err := models.Dbmap.Select(&epics, "SELECT * FROM Epic WHERE id IN (SELECT epicid FROM EpicUserMap WHERE userid=?)", user_id)
	return epics, err
}

func PostEpic(c *gin.Context) {
	id := c.Params.ByName("id")
	var epic models.Epic
	c.Bind(&epic)

	if epic.Name != "" {

		if insert, _ := models.Dbmap.Exec(`INSERT INTO Epic (name) VALUES (?)`, epic.Name); insert != nil {
			epic_id, err := insert.LastInsertId()
			if err == nil {
				models.Dbmap.Exec(`INSERT INTO EpicUserMap (userid, epicid) VALUES (?, ?)`, id, epic_id)
				content := &models.Epic{
					Id:   epic_id,
					Name: epic.Name,
				}
				c.JSON(http.StatusCreated, content)
			} else {
				utils.CheckErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
	}
}

/*func UpdateEpic(c *gin.Context) {
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
}*/

// Still need to recursively delete all stories that are members of this epic
func DeleteEpic(c *gin.Context) {
	id := c.Params.ByName("id")

	var epic models.Epic
	err := models.Dbmap.SelectOne(&epic, "SELECT id FROM User WHERE id=?", id)

	if err == nil {
		_, err = models.Dbmap.Delete(&epic)

		if err == nil {
			c.JSON(http.StatusOK, gin.H{"id #" + id: " deleted"})
		} else {
			utils.CheckErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

/*
func addUser(c *gin.Context) {

}
*/
