package main

import (
    "ScrumifyBackend/models"
    "ScrumifyBackend/utils"
    "github.com/gin-gonic/gin"
    "strconv"
    "ScrumifyBackend/server"
)

func main() {
    utils.InitializeConfig()
    models.InitializeDb()
    if (utils.Conf.ENV == "prod") {
        gin.SetMode(gin.ReleaseMode)
    }

    r := server.RegisterRoutes()
    r.RunTLS(":" + strconv.FormatInt(utils.Conf.PORT, 10), utils.Conf.CERT_PATH, utils.Conf.KEY_PATH)
}
