package models

// Dependencies between modules (dependee depends on dependency)
type ModuleDependencyMap struct {
	Id           int64 `db:"id" json:"id"`
	DependencyID int64 `db:"dependencyid" json:"dependencyid"`
	DependeeID   int64 `db:"dependeeid" json:"dependeeid"`
}

func GetModuleDependencyMap() {

}

func CreateModuleDependencyMap() {

}

func UpdateModuleDependencyMap() {

}

func DeleteModuleDependencyMap() {

}
