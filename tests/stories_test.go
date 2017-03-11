package tests

import (
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/require"
	"net/http"
	"github.com/stretchr/testify/assert"
	"strconv"
	"ScrumifyBackend/models"
	"database/sql"
	"net/http/httptest"
	"time"
)

type StoriesTest struct {
	suite.Suite
	user_id1 string
	user_id2 string
	epic_id1 string
	epic_id2 string
}

func (suite *StoriesTest) SetupTest() {
	cleanDb()
	suite.createTwoEpicsForTwoUsers()
}

func (suite *StoriesTest) TestStoryDoesntExist() {
	trace()

	// User #1 tries to get story
	resp := getRequestResponse("GET", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to update story
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": 1,
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to delete story
	resp = getRequestResponse("DELETE", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/1", "")
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *StoriesTest) TestNonExistentUserCreateStory() {
	trace()
	i, _ := strconv.ParseInt(suite.epic_id2, 10, 64)

	// User #3 tries to create a story under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/" + strconv.FormatInt(i + 1, 10), `{
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestCreateInvalidStory() {
	trace()
	assert := assert.New(suite.T())

	// User #1 tries to create an invalid story under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"stage": 3,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under an unspecified epic
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{"name": "Test story"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"points": -1,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to create an invalid story under epic #1
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"assigned_to": 2,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)
}

func (suite *StoriesTest) TestCreateDeleteStory() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates two stories under epic #1
	resp := createValidStory(suite)
	resp = createValidStory(suite)
	story := unmarshalToStory(resp)
	id := strconv.FormatInt(story.Id, 10)

	// User #1 gets story
	resp = getRequestResponse("GET", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/" + id, "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 deletes story
	resp = getRequestResponse("DELETE", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to get story
	resp = getRequestResponse("GET", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestUpdateInvalidStory() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates story under epic #1
	resp := createValidStory(suite)
	story := unmarshalToStory(resp)
	id := strconv.FormatInt(story.Id, 10)

	// User #1 tries to change story's name to an invalid name
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + id + `,
		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story's stage to an invalid stage
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + id + `,
  		"name": "Test story",
  		"stage": 3,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story's epic_id
	// TODO: Updating a story shouldn't be able to change the epic it is associated with
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + id + `,
		"name": "Test story"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story's points to an invalid value
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + id + `,
  		"name": "Test story",
  		"points": -1,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change story's assignee to one not in epic #1
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + id + `,
  		"name": "Test story",
  		"assigned_to": ` + suite.user_id2 + `,
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// User #1 tries to change an unspecified story
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestUpdateStory() {
	trace()

	// User #1 creates story under epic #1
	resp := createValidStory(suite)
	story := unmarshalToStory(resp)

	// User #1 changes story
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id1, `{
		"id": ` + strconv.FormatInt(story.Id, 10) + `,
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *StoriesTest) TestAccessUnownedStory() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates story under epic #1
	resp := createValidStory(suite)
	story := unmarshalToStory(resp)
	id := strconv.FormatInt(story.Id, 10)

	// User #2 tries to get story
	resp = getRequestResponse("GET", "/api/v1/stories/" + suite.user_id2 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to change story
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id2, `{
		"id": ` + id + `,
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to delete story
	resp = getRequestResponse("DELETE", "/api/v1/stories/" + suite.user_id2 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #1 tries to create a story under epic #2
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id2 + `}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 creates story under epic #2
	resp = getRequestResponse("POST", "/api/v1/stories/" + suite.user_id2, `{
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id2 + `}`)
	assert.Equal(http.StatusCreated, resp.Code)
}

func (suite *StoriesTest) TestMultiownedEpic() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates story under epic #1
	resp := createValidStory(suite)
	story := unmarshalToStory(resp)
	id := strconv.FormatInt(story.Id, 10)

	// User #1 adds user #2 to epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1 + "/" + suite.epic_id1,
		`{"email": "dur2@dur.com"}`)

	// User #2 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + suite.epic_id1, "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 gets story
	resp = getRequestResponse("GET", "/api/v1/stories/" + suite.user_id2 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 changes story
	resp = getRequestResponse("PUT", "/api/v1/stories/" + suite.user_id2, `{
		"id": ` + id + `,
  		"name": "Test story",
  		"epic_id": ` + suite.epic_id1 + `}`)
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 deletes story
	resp = getRequestResponse("DELETE", "/api/v1/stories/" + suite.user_id2 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to get story
	resp = getRequestResponse("GET", "/api/v1/stories/" + suite.user_id1 + "/" + suite.epic_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *StoriesTest) TestCascadeDeletesStory() {
	trace()

	// User #1 creates story under epic #1
	resp := createValidStory(suite)
	story := unmarshalToStory(resp)

	assert.True(suite.T(), storyExists(suite.epic_id1, story.Id))

	// User #1 deletes epic #1
	resp = getRequestResponse("DELETE", "/api/v1/epics/" + suite.user_id1 + "/" + suite.epic_id1, "")

	// User #1 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/" + suite.epic_id1, "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	time.Sleep(10 * time.Millisecond)
	assert.False(suite.T(), storyExists(suite.epic_id1, story.Id))
}

func storyExists(epic_id string, story_id int64) bool {
	var story models.Story
	if err := models.Dbmap.SelectOne(&story, "SELECT * FROM Story WHERE epic_id=? AND id=?", epic_id, story_id);
			err == sql.ErrNoRows {
		return false
	}
	return true
}

func (suite *StoriesTest)createTwoEpicsForTwoUsers() {
	require := require.New(suite.T())

	// Creates user #1
	resp := getRequestResponse("POST", "/api/v1/users", validUser)
	require.Equal(http.StatusCreated, resp.Code)
	user := unmarshalToUser(resp)
	suite.user_id1 = strconv.FormatInt(user.Id, 10)

	// Creates user #2
	resp = getRequestResponse("POST", "/api/v1/users", validUser2)
	require.Equal(http.StatusCreated, resp.Code)
	user = unmarshalToUser(resp)
	suite.user_id2 = strconv.FormatInt(user.Id, 10)

	// User #1 creates epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	suite.epic_id1 = strconv.FormatInt(epic.Id, 10)

	// User #2 creates epic #2
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id2, `{"name": "Test epic2"}`)
	require.Equal(http.StatusCreated, resp.Code)
	epic = unmarshalToEpic(resp)
	suite.epic_id2 = strconv.FormatInt(epic.Id, 10)

	// User #1 tries to get epic #2
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/" + suite.epic_id2, "")
	require.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + suite.epic_id1, "")
	require.Equal(http.StatusUnauthorized, resp.Code)
}

func createValidStory(suite *StoriesTest) *httptest.ResponseRecorder {
	// User #1 creates story under epic #1
	resp := getRequestResponse("POST", "/api/v1/stories/" + suite.user_id1, `{
  		"name": "Test story",
  		"stage": 1,
  		"points": 2,
  		"assigned_to": ` + suite.user_id1 + `,
  		"epic_id": ` + suite.epic_id1 + `,
  		"dependencies": {
    		"Valid": true}}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	return resp
}
