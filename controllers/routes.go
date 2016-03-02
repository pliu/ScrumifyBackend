package controllers

import (
	"TodoBackend/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	usersv1 := r.Group("api/v1")
	{
		usersv1.GET("/users/:id", GetUser)
		// curl -i http://localhost:8080/api/v1/users/1

		usersv1.POST("/users", PostUser)
		// curl -i -X POST -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashedpw\": \"abc\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users

		usersv1.PUT("/users/:id", UpdateUser)
		// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashedpw\": \"cba\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users/1

		usersv1.DELETE("/users/:id", DeleteUser)
		// curl -i -X DELETE http://localhost:8080/api/v1/users/1

		usersv1.GET("/epics/:id", GetEpics)
		// curl -i http://localhost:8080/api/v1/epics/1

		usersv1.POST("/epics/:id", PostEpic)
		// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Test epic\" }" http://localhost:8080/api/v1/epics/1

		usersv1.DELETE("/epics/:id", DeleteEpic)
		// curl -i -X DELETE http://localhost:8080/api/v1/epics/1
	}

	adminv1 := r.Group("admin/v1", gin.BasicAuth(gin.Accounts{"admin": utils.Conf.ADMIN_PASSWORD}))
	{
		adminv1.GET("/users", GetUsers)
		// curl -i http://localhost:8080/api/v1/users
	}

	return r
}
