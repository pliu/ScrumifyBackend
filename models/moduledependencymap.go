package models

import "gopkg.in/gorp.v2"

// Dependencies between modules (dependee depends on dependency)
type ModuleDependencyMap struct {
	DependencyId int64 `db:"dependencyid" json:"dependencyid"`
	DependeeId   int64 `db:"dependeeid" json:"dependeeid"`
}

func SetModuleDependencyMapProperties(table *gorp.TableMap) {
	table.AddIndex("DependeeIdIndex", "Hash", []string{"DependeeId"})
	table.SetUniqueTogether("DependencyId", "DependeeId")
	table.ColMap("DependencyId").SetNotNull(true)
	table.ColMap("DependeeId").SetNotNull(true)
}

func GetModuleDependencyMap() {

}

func CreateModuleDependencyMap() {

}

func UpdateModuleDependencyMap() {

}

func DeleteModuleDependencyMap() {

}
