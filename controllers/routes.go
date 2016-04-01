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
		// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"username\": \"Updated\", \"hashedpw\": \"cba\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users/1

		usersv1.DELETE("/users/:id", DeleteUser)
		// curl -i -X DELETE http://localhost:8080/api/v1/users/1

		usersv1.GET("/users/:id/epics", GetEpics)
		// curl -i http://localhost:8080/api/v1/users/1/epics

		usersv1.POST("/users/:id/epics", PostEpic)
		// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Test epic\" }" http://localhost:8080/api/v1/users/1/epics

		usersv1.PUT("/users/:id/epics/:epicid", UpdateEpic)
		// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"name\": \"New epic\" }" http://localhost:8080/api/v1/users/1/epics/2

		usersv1.DELETE("/users/:id/epics/:epicid", DeleteEpic)
		// curl -i -X DELETE http://localhost:8080/api/v1/users/1/epics/1

		usersv1.POST("/users/:id/epics/:epicid", AddUserToEpic)
		// curl -i -X POST -H "Content-Type: application/json" -d "{ \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users/1/epics/2

		//usersv1.GET("/users/:id/epics/:epicid/modules", GetModules)

		//usersv1.POST("/users/:id/epics/:epicid/modules", PostModule)

		//usersv1.PUT("/modules/:id", UpdateModule)

		//usersv1.DELETE("/modules/:id", DeleteModule)

		//usersv1.GET("/stories/:id", GetStories)

		//usersv1.POST("/stories/:id", PostModule)

		//usersv1.PUT("/stories/:id", UpdateModule)

		//usersv1.DELETE("/stories/:id", DeleteModule)
	}

	adminv1 := r.Group("admin/v1", gin.BasicAuth(gin.Accounts{"admin": utils.Conf.ADMIN_PASSWORD}))
	{
		adminv1.GET("/users", GetUsers)
		// curl -i http://localhost:8080/admin/v1/users
	}

	return r
}
