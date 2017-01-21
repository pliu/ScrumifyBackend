package tests

import (
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"time"
)

type EpicsTest struct {
	suite.Suite
}

func (suite *EpicsTest) SetupTest() {
	models.Dbmap.TruncateTables()
	createTwoUsers()

	// Gets user #1
	resp := getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Gets user #2
	resp = getRequestResponse("GET", "/api/v1/users/2", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *EpicsTest) TestEpicDoesntExist() {
	// User #1 tries to get epic #1
	resp := getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to change epic #1
	resp = getRequestResponse("PUT", "/api/v1/epics/1", `{
		"id": 1,
		"name": "Test epic2"}`)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to delete epic #1
	resp = getRequestResponse("DELETE", "/api/v1/epics/1/1", "")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestCreateInvalidEpic() {
	// Tries to create an invalid epic
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{}`)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *EpicsTest) TestCreateDeleteEpic() {
	assert := assert.New(suite.T())

	// User #1 tries to get epic #1
	resp := getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 creates epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)
	assert.Equal(http.StatusCreated, resp.Code)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 deletes epic #1
	resp = getRequestResponse("DELETE", "/api/v1/epics/1/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestUpdateInvalidEpic() {
	// User #1 creates epic #1
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 tries to change unspecified epic's name to an invalid name
	resp = getRequestResponse("PUT", "/api/v1/epics/1", `{}`)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)

	// User #1 tries to change epic #1's name to an invalid name
	resp = getRequestResponse("PUT", "/api/v1/epics/1", `{"id": 1}`)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *EpicsTest) TestUpdateEpic() {
	// User #1 creates epic #1
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #1 changes epic #1's name
	resp = getRequestResponse("PUT", "/api/v1/epics/1", `{
		"id": 1,
		"name": "Test epic2"}`)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *EpicsTest) TestAccessUnownedEpic() {
	assert := assert.New(suite.T())

	// User #1 creates epic #1
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #2 adds user #2 to epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/2/1", `{"email": "dur2@dur.com"}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 tries to change epic #1's name
	resp = getRequestResponse("PUT", "/api/v1/epics/2", `{
		"id": 1,
		"name": "Test epic2"}`)
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #1 tries to delete epic #1
	resp = getRequestResponse("DELETE", "/api/v1/epics/2/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *EpicsTest) TestAddUserToEpic() {
	assert := assert.New(suite.T())

	// User #1 creates epic #1
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #1 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// User #2 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 adds user #2 to epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/1/1", `{"email": "dur2@dur.com"}`)
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #2 changes epic #1's name
	resp = getRequestResponse("PUT", "/api/v1/epics/2", `{
		"id": 1,
		"name": "Test epic2"}`)
	assert.Equal(http.StatusOK, resp.Code)
}

func (suite *EpicsTest) TestMultiownedDeleteEpic() {
	assert := assert.New(suite.T())

	// User #1 creates epic #1
	resp := getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)

	// User #1 adds user #2 to epic #1
	resp = getRequestResponse("POST", "/api/v1/epics/1/1", `{"email": "dur2@dur.com"}`)

	// User #2 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Delete user #1
	resp = getRequestResponse("DELETE", "/api/v1/users/1", "")

	// Tries to get user #1
	resp = getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// User #1 tries to get epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/1/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)

	// User #2 gets epic #1
	resp = getRequestResponse("GET", "/api/v1/epics/2/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	assert.True(epicExists(1))

	// Delete user #2
	resp = getRequestResponse("DELETE", "/api/v1/users/2", "")

	// Tries to get user #2
	resp = getRequestResponse("GET", "/api/v1/users/2", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	time.Sleep(100 * time.Millisecond)
	assert.False(epicExists(1))
}

func epicExists(epic_id int) bool {
	var epic models.Epic
	if err := models.Dbmap.SelectOne(&epic, "SELECT * FROM Epic WHERE id=?", epic_id); err == sql.ErrNoRows {
		return false
	}
	return true
}
