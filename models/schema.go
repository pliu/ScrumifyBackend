package models

type User struct {
	Id        int64   `db:"id" json:"id"`
	Username  string  `db:"username" json:"username"`
	HashedPW  string  `db:"hashedpw" json:"hashedpw"`
	Email     string  `db:"email" json:"email"`
}

type MemberUser struct {
	Id      int64  `db:"id" json:"id"`
	UserID  int64  `db:"userid" json:"userid"`
	EpicID  int64  `db:"epicid" json:"epicid"`
}

type MemberStory struct {
	Id       int64  `db:"id" json:"id"`
	StoryID  int64  `db:"storyid" json:"storyid"`
	EpicID   int64  `db:"epicid" json:"epicid"`
}

type Story struct {
	Id           int64   `db:"id" json:"id"`
	Name         string  `db:"name" json:"name"`
	Description  string  `db:"description" json:"description"`
	Index        int64   `db:"index" json:"index"`
	Points       int64   `db:"points" json:"points"`
}

type Epic struct {
	Id    int64   `db:"id" json:"id"`
	Name  string  `db:"name" json:"name"`
}

type Session struct {
	SessionID  int64  `db:"sessionid" json:"sessionid"`
	UserID     int64  `db:"userid" json:"userid"`
}

type Device struct {
	Id    int64   `db:"id" json:"id"`
	Name  string  `db:"name" json:"name"`
}

type MemberDevice struct {
	Id        int64  `db:"id" json:"id"`
	DeviceID  int64  `db:"deviceid" json:"deviceid"`
	UserID    int64  `db:"userid" json:"userid"`
}