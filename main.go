package main

import (
"github.com/gin-gonic/gin"
"strconv"
"database/sql"
_ "github.com/go-sql-driver/mysql"
"gopkg.in/gorp.v1"
"log"
"math/rand"
"TodoBackend/models"
)

func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", GetUsers)           // curl -i http://localhost:8080/api/v1/users
		v1.GET("/users/:id", GetUser)        // curl -i http://localhost:8080/api/v1/users/1
		v1.POST("/users", PostUser)          // curl -i -X POST -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashedpw\": \"abc\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users
		v1.PUT("/users/:id", UpdateUser)     // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashedpw\": \"cba\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users/1
		v1.DELETE("/users/:id", DeleteUser)  // curl -i -X DELETE http://localhost:8080/api/v1/users/1
		v1.POST("/auth", Login)
		v1.DELETE("/auth", Logoff)
	}

	r.Run(":8080")
}

func initDB() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:blahblah@/" + DATABASE_NAME)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTable(models.User{}).SetKeys(true, "Id")
	dbmap.AddTable(models.MemberUser{}).SetKeys(true, "Id")
	dbmap.AddTable(models.MemberStory{}).SetKeys(true, "Id")
	dbmap.AddTable(models.Story{}).SetKeys(true, "Id")
	dbmap.AddTable(models.Epic{}).SetKeys(true, "Id")
	dbmap.AddTable(models.Session{}).SetKeys(false, "SessionID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

var dbmap = initDB()

func Login(c *gin.Context) {
	tentativeSession := rand.Int63()
	var user, cmp models.User
	c.Bind(&user)
	if user.Username != "" && user.HashedPW != ""{
		err := dbmap.SelectOne(&cmp, "SELECT * FROM User WHERE username=?", user.Username)
		if err == nil && cmp.HashedPW == user.HashedPW {
				dbmap.Exec(`INSERT INTO Session (sessionid, userid) VALUES (?, ?)`, tentativeSession, cmp.Id);
				c.JSON(201, gin.H{"user": cmp, "token": tentativeSession})
		} else {
			c.JSON(404, gin.H{"error": "user not found"})
		}

	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}

// curl -i -X POST -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashedpw\": \"cba\" }" http://localhost:8080/api/v1/auth
}

func Logoff(c *gin.Context) {
	var user, session models.Session
	c.Bind(&user)
	err := dbmap.SelectOne(&session, "SELECT * FROM Session WHERE sessionid=?", user.SessionID)

	if err == nil {
		_, err = dbmap.Delete(&session)

		if err == nil {
			c.JSON(200, session)
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, session)
	}
}

func AddEpic(c *gin.Context) {
	
}

func DeleteEpic(c *gin.Context) {

}

func AddStory(c *gin.Context) {

}

func DeleteStory(c *gin.Context) {

}

func EditStory(c *gin.Context) {
	
}