package main

import (
    "ScrumifyBackend/controllers"
    "ScrumifyBackend/models"
    "ScrumifyBackend/utils"
    "github.com/gin-gonic/gin"
    "strconv"
)

func main() {
    utils.InitializeConfig()
    models.InitializeDb()
    if (utils.Conf.ENV == "prod") {
        gin.SetMode(gin.ReleaseMode)
    }

    r := controllers.RegisterRoutes()
    r.Run(":" + strconv.FormatInt(utils.Conf.PORT, 10))
}
