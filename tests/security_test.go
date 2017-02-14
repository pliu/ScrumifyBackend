package tests

import (
    "github.com/stretchr/testify/suite"
    "ScrumifyBackend/models"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
    "net/http"
	"strconv"
)

var maliciousUser string = `{
		"username": "; DROP TABLE Story",
		"hashed_pw": "ableh",
		"email": "dur@dur.com"}`
var maliciousUser2 string = `{
		"username": "; DROP TABLE Story",
		"hashed_pw": "ableh",
		"email": "; DROP TABLE Story"}`

type SecurityTest struct {
    suite.Suite
}

func (suite *SecurityTest) SetupTest() {
	cleanDb()
}

func (suite *SecurityTest) TestMissingTableDetectable() {
	trace()

    require.True(suite.T(), storyTableExists())

    models.Dbmap.Exec("DROP TABLE Story")
    assert.False(suite.T(), storyTableExists())

	models.InitializeDb()
}

func (suite *SecurityTest) TestSQLInjection() {
	trace()

    assert := assert.New(suite.T())
    require := require.New(suite.T())

    require.True(storyTableExists())

    resp := getRequestResponse("POST", "/api/v1/users", maliciousUser)
    require.Equal(http.StatusCreated, resp.Code)
	user := unmarshalToUser(resp)
    assert.True(storyTableExists())

    resp = getRequestResponse("PUT", "/api/v1/users/%3BDROP%20TABLE%20Story", maliciousUser2)
    assert.Equal(http.StatusBadRequest, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("PUT", "/api/v1/users/" + strconv.FormatInt(user.Id, 10), maliciousUser2)
    require.Equal(http.StatusOK, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("GET", "/api/v1/users/%3BDROP%20TABLE%20Story", "")
    assert.Equal(http.StatusUnauthorized, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("DELETE", "/api/v1/users/%3BDROP%20TABLE%20Story", "")
    assert.Equal(http.StatusOK, resp.Code)
    assert.True(storyTableExists())
}

func (suite *SecurityTest) TestDoesntLeakHashedPw() {
	trace()
	require := require.New(suite.T())
	assert := assert.New(suite.T())

	// Creates user #1
	resp := getRequestResponse("POST", "/api/v1/users", validUser)
	require.Equal(http.StatusCreated, resp.Code)
	user := unmarshalToUser(resp)
	assert.Equal("", user.HashedPw)
	user_id1 := strconv.FormatInt(user.Id, 10)

	// Creates user #2
	resp = getRequestResponse("POST", "/api/v1/users", validUser2)
	require.Equal(http.StatusCreated, resp.Code)
	user = unmarshalToUser(resp)
	user_id2 := strconv.FormatInt(user.Id, 10)

	// User creates epic
	resp = getRequestResponse("POST", "/api/v1/epics/" + user_id1, `{"name": "Test epic"}`)
	require.Equal(http.StatusCreated, resp.Code)
	epic := unmarshalToEpic(resp)
	epic_id := strconv.FormatInt(epic.Id, 10)
	for _, u := range epic.Members {
		assert.Equal("", u.HashedPw)
	}

	// User #1 adds user #2 to epic
	resp = getRequestResponse("POST", "/api/v1/epics/" + user_id1 + "/" + epic_id, `{"email": "dur2@dur.com"}`)

	// User #2 gets epic
	resp = getRequestResponse("GET", "/api/v1/epics/" + user_id2 + "/" + epic_id, "")
	require.Equal(http.StatusOK, resp.Code)
	epic = unmarshalToEpic(resp)
	for _, u := range epic.Members {
		assert.Equal("", u.HashedPw)
	}

	// Gets user #1
	resp = getRequestResponse("GET", "/api/v1/users/" + user_id1, "")
	require.Equal(http.StatusOK, resp.Code)
	user = unmarshalToUser(resp)
	assert.Equal("", user.HashedPw)

	// User #1 updates user #1
	resp = getRequestResponse("PUT", "/api/v1/users/" + user_id1, `{
		"username": "test3",
  		"hashed_pw": "ableh3",
  		"email": "dur3@dur.com"}`)
	require.Equal(http.StatusOK, resp.Code)
	user = unmarshalToUser(resp)
	assert.Equal("", user.HashedPw)

	// User #1 updates epic
	resp = getRequestResponse("PUT", "/api/v1/epics/" + user_id1, `{
		"id": ` + epic_id + `,
		"name": "Test epic2"}`)
	require.Equal(http.StatusOK, resp.Code)
	epic = unmarshalToEpic(resp)
	for _, u := range epic.Members {
		assert.Equal("", u.HashedPw)
	}
}

func storyTableExists() bool {
	if _, err := models.Dbmap.Exec("SELECT * FROM Story"); err != nil {
		return false
	}
	return true
}
