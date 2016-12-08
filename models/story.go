package models

import (
    "gopkg.in/gorp.v2"
    "ScrumifyBackend/utils"
    "time"
    "database/sql"
    "strconv"
    "ScrumifyBackend/scrumifytypes"
)

type Story struct {
    Id           int64                      `db:"id" json:"id"`
    Name         string                     `db:"name" json:"name"`
    Description  string                     `db:"description" json:"description"`
    DueDate      gorp.NullTime              `db:"due_date" json:"due_date"`
    Points       int64                      `db:"points" json:"points"`
    Stage        int64                      `db:"stage" json:"stage"`
    EpicId       int64                      `db:"epic_id" json:"epic_id"`
    AssignedTo   int64                      `db:"assigned_to" json:"assigned_to"`
    Dependencies scrumifytypes.Dependencies `db:"dependencies" json:"dependencies"`
    CreatedAt    time.Time                  `db:"created_at" json:"-"`
    UpdatedAt    time.Time                  `db:"updated_at" json:"updated_at"`
}

func SetStoryProperties(table *gorp.TableMap) {
    table.SetKeys(true, "Id")

    // InnoDB does not have Hash indices
    table.AddIndex("StoryEpicIdIndex", "Btree", []string{"epic_id"})
    table.AddIndex("StoryCreatedAtIndex", "Btree", []string{"created_at"})
    table.AddIndex("StoryUpdatedAtIndex", "Btree", []string{"updated_at"})
    table.ColMap("Name").SetNotNull(true)
    table.ColMap("Stage").SetNotNull(true)
    table.ColMap("EpicId").SetNotNull(true)
    table.ColMap("CreatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP")
    table.ColMap("UpdatedAt").SetNotNull(true).SetDefaultStatement("DEFAULT CURRENT_TIMESTAMP ON UPDATE " +
        "CURRENT_TIMESTAMP")
}

func GetStory(story_id string) (Story, error) {
    var story Story
    if err := Dbmap.SelectOne(&story, "SELECT * FROM Story WHERE id=?", story_id); err == nil {
        return story, err
    } else if err == sql.ErrNoRows {
        return Story{}, utils.StoryDoesntExist
    } else {
        utils.PrintErr(err, "GetStory: Failed to select story " + story_id)
        return Story{}, err
    }
}

func CreateUpdateStory(story Story, update bool) (Story, error) {
    trans, err := Dbmap.Begin()
    if err != nil {
        utils.PrintErr(err, "CreateStory: Failed to begin transaction")
        return Story{}, err
    }

    if update {
        _, err = trans.Update(&story)
    } else {
        err = trans.Insert(&story)
    }
    if err != nil {
        trans.Rollback()
        utils.PrintErr(err, "CreateUpdateStory: Failed to insert/update story " + strconv.FormatInt(story.Id, 10))
        return Story{}, err
    }
    var check Story
    if err = trans.SelectOne(&check, "SELECT * FROM Story WHERE id=?", story.Id); err == nil {
        return check, trans.Commit()
    } else {
        return story, trans.Commit()
    }
}

func DeleteStory(story Story) error {
    _, err := Dbmap.Delete(&story)
    utils.PrintErr(err, "DeleteStory: Failed to delete story " + strconv.FormatInt(story.Id, 10))
    return err
}

func (story Story)IsValid() bool {
    return story.Name != "" && (story.Stage == 0 || story.Stage == 1 || story.Stage == 2) && story.EpicId >= 1 &&
        story.Points >= 0 && story.Dependencies.IsValid()
}
