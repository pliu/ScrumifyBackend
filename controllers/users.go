package controllers

func GetUsers(c *gin.Context) {
	var users []models.User
	_, err := dbmap.Select(&users, "SELECT * FROM User")

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &models.User{
			Id: user_id,
			Username: user.Username,
			HashedPW: user.HashedPW,
			Email: user.Email,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}

func PostUser(c *gin.Context) {
	var user models.User
	c.Bind(&user)

	if user.Username != "" && user.HashedPW != "" && user.Email != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO User (username, hashedpw, email) VALUES (?, ?, ?)`, user.Username, user.HashedPW, user.Email); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				content := &models.User{
					Id: user_id,
					Username: user.Username,
					HashedPW: user.HashedPW,
					Email: user.Email,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		var json models.User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := models.User{
			Id: user_id,
			Username: json.Username,
			HashedPW: json.HashedPW,
			Email: json.Email,
		}

		if user.Username != "" && user.HashedPW != "" && user.Email != ""{
			_, err = dbmap.Update(&user)

			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(422, gin.H{"error": "fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	var user models.User
	err := dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}