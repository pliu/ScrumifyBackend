package tests

import (
    "github.com/stretchr/testify/suite"
    "ScrumifyBackend/models"
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
    "net/http"
)

var maliciousUser string = `{
		"username": "; DROP TABLE Story",
		"hashed_pw": "ableh",
		"email": "dur@dur.com"}`
var maliciousUser2 string = `{
		"username": "; DROP TABLE Story",
		"hashed_pw": "ableh",
		"email": "; DROP TABLE Story"}`

func TestSuiteSecurityTest(t *testing.T) {
    suite.Run(t, new(SecurityTest))
}

type SecurityTest struct {
    suite.Suite
}

func (suite *SecurityTest) SetupTest() {
    models.InitializeDb()
}

func (suite *SecurityTest) TestMissingTableDetectable() {
    require.True(suite.T(), storyTableExists())

    models.Dbmap.Exec("DROP TABLE Story")
    require.False(suite.T(), storyTableExists())
}

func (suite *SecurityTest) TestSQLInjection() {
    assert := assert.New(suite.T())
    require := require.New(suite.T())

    require.True(storyTableExists())

    resp := getRequestResponse("POST", "/api/v1/users", maliciousUser)
    require.Equal(http.StatusCreated, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("PUT", "/api/v1/users/%3BDROP%20TABLE%20Story", maliciousUser2)
    assert.Equal(http.StatusBadRequest, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("PUT", "/api/v1/users/1", maliciousUser2)
    require.Equal(http.StatusOK, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("GET", "/api/v1/users/%3BDROP%20TABLE%20Story", "")
    assert.Equal(http.StatusUnauthorized, resp.Code)
    assert.True(storyTableExists())

    resp = getRequestResponse("DELETE", "/api/v1/users/%3BDROP%20TABLE%20Story", "")
    assert.Equal(http.StatusOK, resp.Code)
    assert.True(storyTableExists())
}

func storyTableExists() bool {
    if _, err := models.Dbmap.Exec("SELECT * FROM Story"); err != nil {
        return false
    }
    return true
}