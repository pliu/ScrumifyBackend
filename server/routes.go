package server

import (
    "ScrumifyBackend/utils"
    "github.com/gin-gonic/gin"
    "ScrumifyBackend/controllers"
)

func RegisterRoutes() *gin.Engine {
    r := gin.Default()

    usersv1 := r.Group("api/v1")
    {
        usersv1.GET("/users/:id", controllers.GetUser)
        // curl -i http://localhost:8080/api/v1/users/1

        usersv1.POST("/users", controllers.PostUser)
        // curl -i -X POST -H "Content-Type: application/json" -d "{ \"username\": \"Test\", \"hashed_pw\": \"abc\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users

        usersv1.PUT("/users/:id", controllers.UpdateUser)
        // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"username\": \"Updated\", \"hashed_pw\": \"cba\", \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/users/1

        usersv1.DELETE("/users/:id", controllers.DeleteUser)
        // curl -i -X DELETE http://localhost:8080/api/v1/users/1

        usersv1.GET("/epics/:id/:epicid", controllers.GetEpic)
        // curl -i http://localhost:8080/api/v1/epics/1

        usersv1.POST("/epics/:id", controllers.PostEpic)
        // curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Test epic\" }" http://localhost:8080/api/v1/epics/1

        usersv1.PUT("/epics/:id", controllers.UpdateEpic)
        // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"id\": 1, \"name\": \"New epic\" }" http://localhost:8080/api/v1/epics/1

        usersv1.DELETE("/epics/:id/:epicid", controllers.DeleteEpic)
        // curl -i -X DELETE http://localhost:8080/api/v1/epics/1/1

        usersv1.POST("/epics/:id/:epicid", controllers.AddUserToEpic)
        // curl -i -X POST -H "Content-Type: application/json" -d "{ \"email\": \"test@test.com\" }" http://localhost:8080/api/v1/epics/1/2

        usersv1.GET("/stories/:id/:storyid", controllers.GetStory)
        // curl -i http://localhost:8080/api/v1/stories/1/1

        usersv1.POST("/stories/:id", controllers.PostStory)
        // curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Test story\", \"stage\": 1, \"epic_id\": 1 }" http://localhost:8080/api/v1/stories/1

        usersv1.PUT("/stories/:id", controllers.UpdateStory)
        // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"id\": 1, \"name\": \"Test story\", \"stage\": 2, \"epic_id\": 1 }" http://localhost:8080/api/v1/stories/1

        usersv1.DELETE("/stories/:id/:storyid", controllers.DeleteStory)
        // curl -i -X DELETE http://localhost:8080/api/v1/stories/1/1
    }

    if (utils.Conf.ADMIN_USERNAME != "" && utils.Conf.ADMIN_PASSWORD != "") {
        adminv1 := r.Group("admin/v1", gin.BasicAuth(gin.Accounts{utils.Conf.ADMIN_USERNAME: utils.Conf.ADMIN_PASSWORD}))
        {
            adminv1.GET("/users", controllers.GetUsers)
            // curl -i http://localhost:8080/admin/v1/users
        }
    }

    return r
}
