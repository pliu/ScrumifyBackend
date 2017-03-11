package models

import (
	"strings"
	"gopkg.in/gorp.v2"
	"ScrumifyBackend/utils"
	"database/sql"
	"strconv"
	"time"
)

type User struct {
	Id        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	HashedPw  string    `db:"hashed_pw" json:"hashed_pw"`
	Email     string    `db:"email" json:"email"`
	Epics     []Epic    `db:"-" json:"epics"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func SetUserProperties(table *gorp.TableMap) {
	table.SetKeys(true, "Id")

	// InnoDB does not have Hash indices
	table.AddIndex("UserCreatedAtIndex", "Btree", []string{"created_at"})
	table.AddIndex("UserUpdatedAtIndex", "Btree", []string{"updated_at"})
	table.ColMap("Username").SetNotNull(true)
	table.ColMap("HashedPw").SetNotNull(true)
	table.ColMap("Email").SetUnique(true).SetNotNull(true)
	table.ColMap("CreatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP")
	table.ColMap("UpdatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP ON UPDATE " +
			"CURRENT_TIMESTAMP")
}

func GetUsers() ([]User, error) {
	var users []User
	_, err := Dbmap.Select(&users, "SELECT * FROM User")
	utils.PrintErr(err, "GetUsers")
	return users, err
}

func GetUser(user_id string) (User, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "GetUser: Failed to begin transaction")
		return User{}, err
	}

	var user User
	if err = trans.SelectOne(&user, "SELECT * FROM User WHERE id=?", user_id); err != nil {
		if err == sql.ErrNoRows {
			trans.Rollback()
			return User{}, utils.UserDoesntExist
		} else {
			trans.Rollback()
			utils.PrintErr(err, "GetUser: Failed to select user " + user_id)
			return User{}, err
		}
	}

	if _, err = trans.Select(&user.Epics, "SELECT * FROM Epic WHERE id IN (SELECT epic_id FROM EpicUserMap WHERE " +
			"user_id=?)", user_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetUser: Failed to select epics for user " + user_id)
		return User{}, err
	}

	var epicIds []int64
	for _, epic := range user.Epics {
		epicIds = append(epicIds, epic.Id)
	}

	if len(epicIds) == 0 {
		return user, trans.Commit()
	}

	epicIdsString := utils.ConvertInt64ArrayToString(epicIds)
	var stories []Story
	if _, err = trans.Select(&stories, "SELECT * FROM Story WHERE epic_id IN " + epicIdsString); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetUser: Failed to select stories for user " + user_id)
		return User{}, err
	}
	var mappings []EpicUserMap
	if _, err = trans.Select(&mappings, "SELECT * FROM EpicUserMap WHERE epic_id IN " + epicIdsString);
			err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetUser: Failed to select mappings for user " + user_id)
		return User{}, err
	}
	var members []User
	if _, err = trans.Select(&members, "SELECT * FROM User WHERE id IN " + epicIdsString); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "GetUser: Failed to select members for user " + user_id)
		return User{}, err
	}
	err = trans.Commit()

	epicStoryMap := make(map[int64][]Story)
	for _, story := range stories {
		if _, ok := epicStoryMap[story.EpicId]; !ok {
			epicStoryMap[story.EpicId] = make([]Story, 0)
		}
		epicStoryMap[story.EpicId] = append(epicStoryMap[story.EpicId], story)
	}
	for i, epic := range user.Epics {
		user.Epics[i].Stories = epicStoryMap[epic.Id]
	}
	idMemberMap := make(map[int64]User)
	for _, member := range members {
		idMemberMap[member.Id] = member
	}
	epicMemberMap := make(map[int64][]User)
	for _, mapping := range mappings {
		if _, ok := epicMemberMap[mapping.EpicId]; !ok {
			epicMemberMap[mapping.EpicId] = make([]User, 0)
		}
		epicMemberMap[mapping.EpicId] = append(epicMemberMap[mapping.EpicId], idMemberMap[mapping.UserId])
	}
	for i, epic := range user.Epics {
		user.Epics[i].Members = epicMemberMap[epic.Id]
	}

	return user, err
}

func CreateUpdateUser(user User, update bool) (User, error) {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "CreateUpdateUser: Failed to begin transaction")
		return User{}, err
	}
	user.Email = strings.ToLower(user.Email)

	if update {
		_, err = trans.Update(&user)
	} else {
		err = trans.Insert(&user)
	}
	if err != nil {
		trans.Rollback()
		if utils.ParseSQLError(err) == utils.SqlDuplicate {
			return User{}, utils.EmailExists
		}
		utils.PrintErr(err, "CreateUpdateUser: Failed to insert/update user " + strconv.FormatInt(user.Id, 10))
		return User{}, err
	}
	return user, trans.Commit()
}

func DeleteUser(user_id string) error {
	trans, err := Dbmap.Begin()
	if err != nil {
		utils.PrintErr(err, "DeleteUser: Failed to begin transaction")
		return err
	}

	var mappings []EpicUserMap
	if _, err = trans.Select(&mappings, "SELECT * FROM EpicUserMap WHERE user_id=?", user_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "DeleteUser: Failed to select mappings for user " + user_id)
		return err
	}
	if _, err = trans.Exec("DELETE FROM User WHERE id=?", user_id); err != nil {
		trans.Rollback()
		utils.PrintErr(err, "DeleteUser: Failed to delete user " + user_id)
		return err
	}
	err = trans.Commit()
	for _, mapping := range mappings {
		go removeUnownedEpic(mapping.EpicId)
	}
	return err
}

func (user User)IsValid() bool {
	return user.Username != "" && user.HashedPw != "" && user.Email != ""
}

func (user User)Scrub() User {
	for _, epic := range user.Epics {
		for i, member := range epic.Members {
			epic.Members[i] = member.scrubPassword()
		}
	}
	return user.scrubPassword()
}

func (user User)scrubPassword() User {
	user.HashedPw = ""
	return user
}
