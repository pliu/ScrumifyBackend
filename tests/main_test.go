package tests

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"ScrumifyBackend/utils"
	"ScrumifyBackend/server"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"bytes"
	"net/http"
)

var r *gin.Engine

func init() {
	utils.InitializeConfig()
	utils.Conf.ENV = "test"
	r = server.RegisterRoutes()
}

func TestSuiteMainTest(t *testing.T) {
	suite.Run(t, new(UsersTest))
	suite.Run(t, new(EpicsTest))
	suite.Run(t, new(StoriesTest))
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
