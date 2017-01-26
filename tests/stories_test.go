package tests

import (
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"time"
)

var validStory string = `{
  		"name": "Test story",
  		"stage": 1,
  		"points": 2,
  		"assigned_to": 1,
  		"epic_id": 1,
  		"dependencies": {
    		"Valid": true}}`

type StoriesTest struct {
	suite.Suite
}

func (suite *StoriesTest) SetupTest() {
	models.InitializeDb()
	createTwoUsers()

	// Gets user #1
	resp := getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Gets user #2
	resp = getRequestResponse("GET", "/api/v1/users/2", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 creates epic #1
	getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #2 creates epic #2
	getRequestResponse("POST", "/api/v1/epics/2", `{"name": "Test epic2"}`)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 tries to get epic #2
	resp = getRequestResponse("GET", "/api/v1/epics/1/2", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #2 gets epic #2
	resp = getRequestResponse("GET", "/api/v1/epics/2/2", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestStoryDoesntExist() {
	// User #1 tries to get story #1
	resp := getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to update story #1
	resp = getRequestResponse("PUT", "/api/v1/stories/1", validStory)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to delete story #1
	resp = getRequestResponse("DELETE", "/api/v1/stories/1/1", "")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestNonExistentUserCreateStory() {
	// User #3 tries to create a story under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/3", `{
  		"name": "Test story",
  		"epic_id": 1}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestCreateInvalidStory() {
	assert := assert.New(suite.T())

	// User #1 tries to create an invalid story under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", `{"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"stage": 3,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under an unspecified epic
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{"name": "Test story"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"points": -1,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"assigned_to": 2,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)
}

func (suite *StoriesTest) TestCreateDeleteStory() {
	assert := assert.New(suite.T())

	// User #1 tries to get story #1
	resp := getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 creates story #1 under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/1", validStory)
	assert.Equal(http.StatusCreated, resp.Code)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 deletes story #1
	resp = getRequestResponse("DELETE", "/api/v1/stories/1/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to get story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestUpdateInvalidStory() {
	assert := assert.New(suite.T())

	// User #1 creates story #1 under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", validStory)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 tries to change story #1's name to an invalid name
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story #1's stage to an invalid stage
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
  		"name": "Test story",
  		"stage": 3,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story #1's epic_id
	// TODO: Updating a story shouldn't be able to change the epic it is associated with
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
		"name": "Test story"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story #1's points to an invalid value
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
  		"name": "Test story",
  		"points": -1,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story #1's assignee to one not in epic #1
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
  		"name": "Test story",
  		"assigned_to": 2,
  		"epic_id": 1}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change an unspecified story
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"epic_id": 1}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestUpdateStory() {
	// User #1 creates story #1 under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", validStory)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 changes story #1
	resp = getRequestResponse("PUT", "/api/v1/stories/1", `{
		"id": 1,
  		"name": "Test story",
  		"epic_id": 1}`)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *StoriesTest) TestAccessUnownedStory() {
	assert := assert.New(suite.T())

	// User #1 creates story #1 under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", validStory)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 tries to get story #1
	resp = getRequestResponse("GET", "/api/v1/stories/2/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to change story #1
	resp = getRequestResponse("PUT", "/api/v1/stories/2", `{
		"id": 1,
  		"name": "Test story",
  		"epic_id": 1}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to delete story #1
	resp = getRequestResponse("DELETE", "/api/v1/stories/2/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #1 tries to create a story under epic #2
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"epic_id": 2}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestMultiownedEpic() {
	assert := assert.New(suite.T())

	// User #1 creates story #1 under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", validStory)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 adds user #2 to epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/1/1", `{"email": "dur2@dur.com"}`)

	// User #2 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/2/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 changes story #1
	resp = getRequestResponse("PUT", "/api/v1/stories/2", `{
		"id": 1,
  		"name": "Test story",
  		"epic_id": 1}`)
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 deletes story #1
	resp = getRequestResponse("DELETE", "/api/v1/stories/2/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 adds user #1 to epic #2
	resp = getRequestResponse("POST", "/api/v1/epics/2/2", `{"email": "dur@dur.com"}`)

	// User #1 gets epic #2
	resp = getRequestResponse("GET", "/api/v1/epics/1/2", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 tries to create a story under epic #2
	resp = getRequestResponse("POST", "/api/v1/stories/1", `{
  		"name": "Test story",
  		"epic_id": 2}`)
	assert.Equal(http.StatusCreated, resp.Code)
}

func (suite *StoriesTest) TestCascadeDeletesStory() {
	// User #1 creates story #1 under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/1", validStory)

	// User #1 gets story #1
	resp = getRequestResponse("GET", "/api/v1/stories/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	assert.True(suite.T(), storyExists(1))

	// User #1 deletes epic #1
	resp = getRequestResponse("DELETE", "/api/v1/epics/1/1", "")

	// User #1 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	time.Sleep(100 * time.Millisecond)
	assert.False(suite.T(), storyExists(1))
}

func storyExists(story_id int) bool {
	var story models.Story
	if err := models.Dbmap.SelectOne(&story, "SELECT * FROM Story WHERE id=?", story_id); err == sql.ErrNoRows {
		return false
	}
	return true
}
