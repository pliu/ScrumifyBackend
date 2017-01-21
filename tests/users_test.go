package tests

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/models"
)

var validUser string = `{
		"username": "test",
		"hashed_pw": "ableh",
		"email": "dur@dur.com"}`
var validUser2 string = `{
		"username": "test2",
  		"hashed_pw": "ableh2",
  		"email": "dur2@dur.com"}`

type UsersTest struct {
	suite.Suite
}

func (suite *UsersTest) SetupTest() {
	models.Dbmap.TruncateTables()
}

func (suite *UsersTest) TestUserDoesntExist() {
	assert := assert.New(suite.T())

	// Tries to get user #1
	resp := getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// Tries to change user #1's name and e-mail
	// TODO: Currently passes, but should fail after authentication is implemented
	resp = getRequestResponse("PUT", "/api/v1/users/1", validUser)
	assert.Equal(http.StatusOK, resp.Code)

	// Tries to delete user #1
	// TODO: Currently passes, but should fail after authentication is implemented
	resp = getRequestResponse("DELETE", "/api/v1/users/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// User #1 tries to create epic
	// TODO: Currently passes, but should fail after authentication is implemented
	resp = getRequestResponse("POST", "/api/v1/epics/1", `{"name": "Test epic"}`)
	assert.Equal(http.StatusCreated, resp.Code)
}

func (suite *UsersTest) TestCreateInvalidUser() {
	assert := assert.New(suite.T())

	// Tries to create an invalid user
	resp := getRequestResponse("POST", "/api/v1/users", `{
		"hashed_pw": "ableh",
		"email": "dur@dur.com"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// Tries to create an invalid user
	resp = getRequestResponse("POST", "/api/v1/users", `{
		"username": "test",
		"email": "dur@dur.com"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// Tries to create an invalid user
	resp = getRequestResponse("POST", "/api/v1/users", `{
		"username": "test",
		"hashed_pw": "ableh",}`)
	assert.Equal(http.StatusBadRequest, resp.Code)
}

func (suite *UsersTest) TestCreateDeleteUser() {
	assert := assert.New(suite.T())

	// Tries to get user #1
	resp := getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusUnauthorized, resp.Code)

	// Creates user #1
	resp = getRequestResponse("POST", "/api/v1/users", validUser)
	assert.Equal(http.StatusCreated, resp.Code)

	// Gets user #1
	resp = getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Delete user #1
	resp = getRequestResponse("DELETE", "/api/v1/users/1", "")
	assert.Equal(http.StatusOK, resp.Code)

	// Tries to get user #1
	resp = getRequestResponse("GET", "/api/v1/users/1", "")
	assert.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *UsersTest) TestCreateUpdateUserDuplicateEmail() {
	createTwoUsers()

	// Gets user #1
	resp := getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Gets user #2
	resp = getRequestResponse("GET", "/api/v1/users/2", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Tries to create a user with the same e-mail
	resp = getRequestResponse("POST", "/api/v1/users", validUser)
	assert.Equal(suite.T(), http.StatusConflict, resp.Code)

	// Tries to change user #2's e-mail to that of user #1
	// TODO: CreateUpdateUser in user.go does not currently distinguish between duplicate e-mails or internal server
	// error
	resp = getRequestResponse("PUT", "/api/v1/users/2", validUser)
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *UsersTest) TestUpdateInvalidUser() {
	assert := assert.New(suite.T())

	// Creates user #1
	resp := getRequestResponse("POST", "/api/v1/users", validUser)

	// Gets user #1
	resp = getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Tries to change user #1's username to an invalid username
	resp = getRequestResponse("PUT", "/api/v1/users/1", `{
  		"hashed_pw": "ableh",
  		"email": "dur@dur.com"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// Tries to change user #1's hashed pw to an invalid hashed pw
	resp = getRequestResponse("PUT", "/api/v1/users/1", `{
  		"username": "test2",
  		"email": "dur@dur.com"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)

	// Tries to change user #1's e-mail to an invalid e-mail
	resp = getRequestResponse("PUT", "/api/v1/users/1", `{
  		"username": "test2",
  		"hashed_pw": "ableh"}`)
	assert.Equal(http.StatusBadRequest, resp.Code)
}

func (suite *UsersTest) TestUpdateUser() {
	// Creates user #1
	resp := getRequestResponse("POST", "/api/v1/users", validUser)

	// Gets user #1
	resp = getRequestResponse("GET", "/api/v1/users/1", "")
	require.Equal(suite.T(), http.StatusOK, resp.Code)

	// Changes user #1's info
	resp = getRequestResponse("PUT", "/api/v1/users/1", validUser2)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	// Changes user #1's e-mail to the same e-mail
	resp = getRequestResponse("PUT", "/api/v1/users/1", validUser)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func createTwoUsers () {
	// Creates user #1
	getRequestResponse("POST", "/api/v1/users", validUser)

	// Creates user #2
	getRequestResponse("POST", "/api/v1/users", validUser2)
}
