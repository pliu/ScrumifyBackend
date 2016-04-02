package models

import (
	"TodoBackend/utils"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

var Dbmap *gorp.DbMap

func init() {
	db, err := sql.Open("mysql", "root:"+utils.Conf.DB_PASSWORD+"@/"+utils.Conf.DB_NAME)
	utils.CheckErr(err, "sql.Open failed")
	Dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	Dbmap.AddTable(User{}).SetKeys(true, "Id")
	Dbmap.AddTable(EpicUserMap{}).SetKeys(true, "Id")
	Dbmap.AddTable(EpicModuleMap{}).SetKeys(false, "ModuleId")
	Dbmap.AddTable(ModuleStoryMap{}).SetKeys(false, "StoryId")
	Dbmap.AddTable(ModuleDependencyMap{}).SetKeys(true, "Id")
	Dbmap.AddTable(Story{}).SetKeys(true, "Id")
	Dbmap.AddTable(Module{}).SetKeys(true, "Id")
	Dbmap.AddTable(Epic{}).SetKeys(true, "Id")
	err = Dbmap.CreateTablesIfNotExists()
	utils.CheckErr(err, "Create table failed")
}
