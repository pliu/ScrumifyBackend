package tests

import (
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"strconv"
	"time"
)

type EpicsTest struct {
	suite.Suite
	user_id1 string
	user_id2 string
}

func (suite *EpicsTest) SetupTest() {
	cleanDb()
	suite.createTwoUsers()
}

func (suite *EpicsTest) TestEpicDoesntExist() {
	trace()

	// User #1 tries to get epic
	resp := getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to change epic
	resp = getRequestResponse("PUT", "/api/v1/epics/" + suite.user_id1, `{
		"id": 1,
		"name": "Test epic2"}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to delete epic
	resp = getRequestResponse("DELETE", "/api/v1/epics/" + suite.user_id1 + "/1", "")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestCreateInvalidEpic() {
	trace()

	// Tries to create an invalid epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{}`)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *EpicsTest) TestCreateDeleteEpic() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	assert.Equal(http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	id := strconv.FormatInt(epic.Id, 10)

	// User #1 gets epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/" + id, "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 deletes epic
	resp = getRequestResponse("DELETE", "/api/v1/epics/" + suite.user_id1 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to get epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestUpdateInvalidEpic() {
	trace()

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)

	// User #1 tries to change unspecified epic
	resp = getRequestResponse("PUT", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to change epic's name to an invalid name
	resp = getRequestResponse("PUT", "/api/v1/epics/"  + suite.user_id1, `{"id": ` +
			strconv.FormatInt(epic.Id, 10) + `}`)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *EpicsTest) TestUpdateEpic() {
	trace()

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)

	// User #1 changes epic's name
	resp = getRequestResponse("PUT", "/api/v1/epics/" + suite.user_id1, `{
		"id": ` + strconv.FormatInt(epic.Id, 10) + `,
		"name": "Test epic2"}`)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *EpicsTest) TestAccessUnownedEpic() {
	trace()

	assert := assert.New(suite.T())

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	id := strconv.FormatInt(epic.Id, 10)

	// User #2 tries to get epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 adds user #2 to epic
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id2 + "/" + id, `{"email": "dur2@dur.com"}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to change epic's name
	resp = getRequestResponse("PUT", "/api/v1/epics/" + suite.user_id2, `{
		"id": ` + id + `,
		"name": "Test epic2"}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to delete epic
	resp = getRequestResponse("DELETE", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestAddUserToEpic() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	id := strconv.FormatInt(epic.Id, 10)

	// User #2 tries to get epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 adds user #2 to epic
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1 + "/" + id, `{"email": "dur2@dur.com"}`)
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 gets epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 changes epic's name
	resp = getRequestResponse("PUT", "/api/v1/epics/" + suite.user_id2, `{
		"id": ` + id + `,
		"name": "Test epic2"}`)
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 deletes epic
	resp = getRequestResponse("DELETE", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)
}

func (suite *EpicsTest) TestMultiownedDeleteEpic() {
	trace()
	assert := assert.New(suite.T())

	// User #1 creates epic
	resp := getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1, `{"name": "Test epic"}`)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	id := strconv.FormatInt(epic.Id, 10)

	// User #1 adds user #2 to epic
	resp = getRequestResponse("POST", "/api/v1/epics/" + suite.user_id1 + "/" + id, `{"email": "dur2@dur.com"}`)

	// User #2 gets epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Delete user #1
	resp = getRequestResponse("DELETE", "/api/v1/users/" + suite.user_id1, "")

	// Tries to get user #1
	resp = getRequestResponse("GET", "/api/v1/users/" + suite.user_id1, "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to get epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id1 + "/" + id, "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 gets epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + suite.user_id2 + "/" + id, "")
	assert.Equal(http.StatusOK, resp.Code)

	assert.True(epicExists(epic.Id))

	// Delete user #2
	resp = getRequestResponse("DELETE", "/api/v1/users/" + suite.user_id2, "")

	// Tries to get user #2
	resp = getRequestResponse("GET", "/api/v1/users/" + suite.user_id2, "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	time.Sleep(10 * time.Millisecond)
	assert.False(epicExists(epic.Id))
}

func epicExists(epic_id int64) bool {
	var epic models.Epic
	if err := models.Dbmap.SelectOne(&epic, "SELECT * FROM Epic WHERE id=?", epic_id); err == sql.ErrNoRows {
		return false
	}
	return true
}

func (suite *EpicsTest)createTwoUsers () {
	// Creates user #1
	resp := getRequestResponse("POST", "/api/v1/users", validUser)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	user := unmarshalToUser(resp)
	suite.user_id1 = strconv.FormatInt(user.Id, 10)

	// Creates user #2
	resp = getRequestResponse("POST", "/api/v1/users", validUser2)
	require.Equal(suite.T(), http.StatusCreated, resp.Code)
	user = unmarshalToUser(resp)
	suite.user_id2 = strconv.FormatInt(user.Id, 10)
}
