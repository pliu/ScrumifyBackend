package models

type Module struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	DueDate string `db:"duedate" json:"duedate"`
	Stage   int64  `db:"stage" json:"stage"`
	Owner   int64  `db:"owner" json:"owner"`
}

func GetModule() {

}

func CreateModule() {

}

func UpdateModule() {

}

func DeleteModule() {

}
