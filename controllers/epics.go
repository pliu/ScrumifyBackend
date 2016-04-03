package controllers

import (
	"TodoBackend/models"
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetEpics(c *gin.Context) {
	id := c.Params.ByName("id")
	var epics []models.Epic
	_, err := models.Dbmap.Select(&epics, "SELECT * FROM Epic WHERE id IN (SELECT epicid FROM EpicUserMap WHERE userid=?)", id)

	if err == nil {
		c.JSON(http.StatusOK, epics)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not access database"})
	}
}

func PostEpic(c *gin.Context) {
	id := c.Params.ByName("id")
	if userExists(id) {
		var epic models.Epic
		c.Bind(&epic)

		if validEpic(epic) {

			if insert, _ := models.Dbmap.Exec(`INSERT INTO Epic (name) VALUES (?)`, epic.Name); insert != nil {
				epic_id, err := insert.LastInsertId()
				if err == nil {
					models.Dbmap.Exec(`INSERT INTO EpicUserMap (userid, epicid) VALUES (?, ?)`, id, epic_id)
					epic.Id = epic_id
					c.JSON(http.StatusCreated, epic)
				} else {
					utils.CheckErr(err, "Insert epic failed")
				}
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
	}
}

func UpdateEpic(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if epicOwnedByUser(id, epic_id) {
		var epic models.Epic
		err := models.Dbmap.SelectOne(&epic, "SELECT * FROM Epic WHERE id=?", epic_id)

		if err == nil {
			var json models.Epic
			c.Bind(&json)

			epic := models.Epic{
				Id:   epic.Id,
				Name: json.Name,
			}

			if validEpic(epic) {
				_, err = models.Dbmap.Update(&epic)

				if err == nil {
					c.JSON(http.StatusOK, epic)
				} else {
					utils.CheckErr(err, "Updated epic failed")
				}

			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Field(s) is(are) empty"})
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Epic not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}

func DeleteEpic(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if epicOwnedByUser(id, epic_id) {
		var mapping models.EpicUserMap
		err := models.Dbmap.SelectOne(&mapping, "SELECT * FROM EpicUserMap WHERE userid=? AND epicid=?", id, epic_id)

		if err == nil {
			_, err = models.Dbmap.Delete(&mapping)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"id #" + epic_id: "Deleted from " + id + "'s list"})
				go removeUnownedEpic(mapping.EpicID)
			} else {
				utils.CheckErr(err, "Delete epic failed")
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Epic not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}

func AddUserToEpic(c *gin.Context) {
	id := c.Params.ByName("id")
	epic_id := c.Params.ByName("epicid")

	if epicOwnedByUser(id, epic_id) {
		var email models.RestEmail
		c.Bind(&email)
		user, err := getUserByEmail(email.Email)

		if err == nil {
			var mapping models.EpicUserMap
			err = models.Dbmap.SelectOne(&mapping, "SELECT * FROM EpicUserMap WHERE userid=? AND epicid=?", user.Id, epic_id)
			if err != nil {
				models.Dbmap.Exec(`INSERT INTO EpicUserMap (userid, epicid) VALUES (?, ?)`, user.Id, epic_id)
				int_epic_id, _ := strconv.ParseInt(epic_id, 10, 64)
				mapping = models.EpicUserMap{
					EpicID: int_epic_id,
					UserID: user.Id,
					Id:     0,
				}
				c.JSON(http.StatusOK, mapping)
			} else {
				c.JSON(http.StatusOK, gin.H{"error": "User already a member of the epic"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Epic not owned by you"})
	}
}
