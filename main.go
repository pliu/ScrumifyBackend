package main

import (
    "TodoBackend/controllers"
    "TodoBackend/models"
    "TodoBackend/utils"
    "github.com/gin-gonic/gin"
    "strconv"
)

/*
Simplifying assumptions (for now):
- users only have one device
- no permissions
- no push

Things to figure out:
- synchronization
*/
func main() {
    utils.InitializeConfig()
    models.InitializeDb()
    if (utils.Conf.ENV == "prod") {
        gin.SetMode(gin.ReleaseMode)
    }

    r := controllers.RegisterRoutes()
    r.Run(":" + strconv.FormatInt(utils.Conf.PORT, 10))
}
