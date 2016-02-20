package main

import (
"github.com/gin-gonic/gin"
"strconv"
"database/sql"
_ "github.com/go-sql-driver/mysql"
"gopkg.in/gorp.v1"
"log"
)

type User struct {
	Id int64 `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	HashedPW string `db:"hashedpw" json:"hashedpw"`
	Email string `db:"email" json:"email"`
}

type MemberUser struct {
	Id int64 `db:"id" json:"id"`
	UserID int64 `db:"userid" json:"userid"`
	EpicID int64 `db:"epicid" json:"epicid"`
}

type MemberStory struct {
	Id int64 `db:"id" json:"id"`
	StoryID int64 `db:"storyid" json:"storyid"`
	EpicID int64 `db:"epicid" json:"epicid"`
}

type Story struct {
	Id int64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type Epic struct {
	Id int64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Session struct {
	Id int64 `db:"id" json:"id"`
	SessionID int64 `db:"sessionid" json:"sessionid"`
	UserID int64 `db:"userid" json:"userid"`
}

func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
		v1.GET("/stories", GetStories)
		v1.GET("/stories/:id", GetStory)
	}

	r.Run(":8080")
}

func initDb(name string) *gorp.DbMap {
	db, err := sql.Open("mysql", "root:blahblah@/myapi")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(User{}, name).SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

var dbmap = initDb("User")

func GetUsers(c *gin.Context) {
	var users []User
	_, err := dbmap.Select(&users, "SELECT * FROM User")

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}

// curl -i http://localhost:8080/api/v1/users
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &User{
			Id: user_id,
			Firstname: user.Firstname,
			Lastname: user.Lastname,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

// curl -i http://localhost:8080/api/v1/users/1
}

func PostUser(c *gin.Context) {
	var user User
	c.Bind(&user)

	if user.Firstname != "" && user.Lastname != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO User (firstname, lastname) VALUES (?, ?)`, user.Firstname, user.Lastname); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				content := &User{
					Id: user_id,
					Firstname: user.Firstname,
					Lastname: user.Lastname,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}

// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		var json User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := User{
			Id: user_id,
			Firstname: json.Firstname,
			Lastname: json.Lastname,
		}

		if user.Firstname != "" && user.Lastname != ""{
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

// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT id FROM user WHERE id=?", id)

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

// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

func GetStories(c *gin.Context) {

}
func GetStory(c *gin.Context) {
	
}