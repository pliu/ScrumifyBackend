package tests

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/utils"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"bytes"
	"net/http"
	"ScrumifyBackend/models"
	"io/ioutil"
	"encoding/json"
	"runtime"
	"fmt"
	"ScrumifyBackend/controllers"
)

var r *gin.Engine

func init() {
	utils.InitializeConfig()
	utils.Conf.ENV = "test"
	models.InitializeDb()
	r = controllers.RegisterRoutes()
}

func TestSuiteMainTest(t *testing.T) {
	suite.Run(t, new(UsersTest))
	suite.Run(t, new(EpicsTest))
	suite.Run(t, new(StoriesTest))
	suite.Run(t, new(SecurityTest))
}

func cleanDb() {
	trans, _ := models.Dbmap.Begin()
	trans.Exec("DELETE FROM User")
	trans.Exec("DELETE FROM Epic")
	trans.Commit()
}

func getRequestResponse(requestType string, endpoint string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(requestType, endpoint, bytes.NewBuffer([]byte(body)))
	} else {
		req, _ = http.NewRequest(requestType, endpoint, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func unmarshalToUser(resp *httptest.ResponseRecorder) models.User {
	var user models.User
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &user)
	return user
}

func unmarshalToEpic(resp *httptest.ResponseRecorder) models.Epic {
	var epic models.Epic
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &epic)
	return epic
}

func unmarshalToStory(resp *httptest.ResponseRecorder) models.Story {
	var story models.Story
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &story)
	return story
}

func trace() {
	pc := make([]uintptr, 10)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("\n%s:%d %s\n", file, line, f.Name())
}
