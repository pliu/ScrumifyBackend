package models

import (
    "ScrumifyBackend/utils"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "gopkg.in/gorp.v2"
    "log"
    "os"
)

var Dbmap *gorp.DbMap

func InitializeDb() {
    var db_name string
    if utils.Conf.ENV == "test" {
        db_name = "todo_test"
    } else if utils.Conf.ENV == "prod" {
        db_name = "todo_prod"
    } else {
        db_name = "todo_dev"
    }

    db, err := sql.Open("mysql", utils.Conf.DB_USERNAME + ":" + utils.Conf.DB_PASSWORD + "@/")
    utils.FatalErr(err, "Connect to database failed")

    if utils.Conf.ENV == "test" {
        _, err = db.Exec("DROP DATABASE IF EXISTS " + db_name)
        utils.FatalErr(err, "Drop database failed")
    }
    _, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + db_name)
    utils.FatalErr(err, "Create database failed")

    db, err = sql.Open("mysql", utils.Conf.DB_USERNAME + ":" + utils.Conf.DB_PASSWORD + "@/" + db_name +
        "?parseTime=true")
    utils.FatalErr(err, "Connect to database failed")
    Dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    if utils.Conf.ENV == "test" || utils.Conf.ENV == "dev" {
        Dbmap.TraceOn("[gorp]", log.New(os.Stdout, "", log.Ltime))
    }

    SetEpicUserMapProperties(Dbmap.AddTable(EpicUserMap{}))
    SetUserProperties(Dbmap.AddTable(User{}))
    SetEpicProperties(Dbmap.AddTable(Epic{}))
    SetStoryProperties(Dbmap.AddTable(Story{}))

    err = Dbmap.CreateTablesIfNotExists()
    utils.FatalErr(err, "Create table failed")

    Dbmap.CreateIndex()
}
