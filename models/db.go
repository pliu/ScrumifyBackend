package models

import (
	"TodoBackend/utils"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

var Dbmap *gorp.DbMap

func init() {
	db, err := sql.Open("mysql", "root:"+utils.Conf.DB_PASSWORD+"@/"+utils.Conf.DB_NAME)
	utils.FatalErr(err, "sql.Open failed")
	Dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	SetEpicUserMapProperties(Dbmap.AddTable(EpicUserMap{}))
	SetModuleDependencyMapProperties(Dbmap.AddTable(ModuleDependencyMap{}))
	SetUserProperties(Dbmap.AddTable(User{}))
	SetEpicProperties(Dbmap.AddTable(Epic{}))
	SetModuleProperties(Dbmap.AddTable(Module{}))
	SetStoryProperties(Dbmap.AddTable(Story{}))
	err = Dbmap.CreateTablesIfNotExists()
	utils.FatalErr(err, "Create table failed")
}
