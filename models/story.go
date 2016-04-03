package models

type Story struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Points      int64  `db:"points" json:"points"`
	Stage       int64  `db:"stage" json:"stage"`
	Owner       int64  `db:"owner" json:"owner"`
}

func GetStory() {

}

func CreateStory() {

}

func UpdateStory() {

}

func DeleteStory() {

}
