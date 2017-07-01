package main

import (
	"encoding/json"
	"fmt"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

type D struct {
	*C
}

func (d *D) AssertTrue(b bool) {
	d.Assert(b, Equals, true)
}

func (d *D) AssertFalse(b bool) {
	d.Assert(b, Equals, false)
}

func (d *D) AssertStatusCode(w *httptest.ResponseRecorder, code int) {
	d.Assert(w.Code, Equals, code)
}

func serve(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	Router().ServeHTTP(w, r)
	return w
}

func runTest(c *C, path string, expectedCode int, expectedResponse Response) {
	d := D{c}

	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", path, nil)
	serve(writer, request)
	d.AssertStatusCode(writer, expectedCode)

	result := writer.Result()
	d.Assert(result.Header["Content-Type"][0], Equals, "application/json")

	var response Response
	_ = json.Unmarshal(writer.Body.Bytes(), &response)

	d.Assert(response, Equals, expectedResponse)
}

func (s *MySuite) TestValidStatusCodes(c *C) {
	var path string
	var expectedResponse Response
	// I know, I know, 1000 is just a random number here
	for i := 1; i <= 1000; i++ {
		if http.StatusText(i) != "" {
			path = fmt.Sprintf("/status/%v", i)
			expectedResponse.Response.Code = i
			expectedResponse.Response.Description = http.StatusText(i)
			runTest(c, path, http.StatusOK, expectedResponse)
		}
	}
}

func (s *MySuite) TestInvalidNumericalStatusCodes(c *C) {
	var path string
	var expectedResponse Response
	for i := 1; i <= 1000; i++ {
		if http.StatusText(i) == "" {
			path = fmt.Sprintf("/status/%v", i)
			expectedResponse.Error = StatusCodeDoesNotExist{code: i}.Error()
			runTest(c, path, http.StatusBadRequest, expectedResponse)
		}
	}
}
