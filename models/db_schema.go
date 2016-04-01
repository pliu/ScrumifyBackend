package models

type User struct {
	Id       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	HashedPW string `db:"hashedpw" json:"hashedpw"`
	Email    string `db:"email" json:"email"`
}

// Which users are part of which epics
type EpicUserMap struct {
	Id     int64 `db:"id" json:"id"`
	UserID int64 `db:"userid" json:"userid"`
	EpicID int64 `db:"epicid" json:"epicid"`
}

// Which modules are part of which epics
type EpicModuleMap struct {
	Id       int64 `db:"id" json:"id"`
	ModuleID int64 `db:"moduleid" json:"moduleid"`
	EpicID   int64 `db:"epicid" json:"epicid"`
}

// Which stories are part of which modules
type ModuleStoryMap struct {
	Id       int64 `db:"id" json:"id"`
	ModuleID int64 `db:"moduleid" json:"moduleid"`
	StoryID  int64 `db:"storyid" json:"storyid"`
}

// Dependencies between modules (dependee depends on dependency)
type ModuleDependencyMap struct {
	Id           int64 `db:"id" json:"id"`
	DependencyID int64 `db:"dependencyid" json:"dependencyid"`
	DependeeID   int64 `db:"dependeeid" json:"dependeeid"`
}

type Story struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Points      int64  `db:"points" json:"points"`
	Stage       int64  `db:"stage" json:"stage"`
}

type Epic struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Module struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	DueDate string `db:"duedate" json:"duedate"`
	Stage   int64  `db:"stage" json:"stage"`
}
