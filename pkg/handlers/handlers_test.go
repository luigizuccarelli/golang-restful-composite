package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-restful-composite/pkg/connectors"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Injected error")
}

// Mock all connections
type Connections struct {
	l    *simple.Logger
	Http *http.Client
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewHttpTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// NewTestConnections - create all mock connections
func NewTestConnections(file string, code int, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &Connections{l: logger, Http: httpclient}
	return conns
}

func (r *Connections) Error(msg string, val ...interface{}) {
	r.l.Error(fmt.Sprintf(msg, val...))
}

func (r *Connections) Info(msg string, val ...interface{}) {
	r.l.Info(fmt.Sprintf(msg, val...))
}

func (r *Connections) Debug(msg string, val ...interface{}) {
	r.l.Debug(fmt.Sprintf(msg, val...))
}

func (r *Connections) Trace(msg string, val ...interface{}) {
	r.l.Trace(fmt.Sprintf(msg, val...))
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestHandlers(t *testing.T) {

	logger := &simple.Logger{Level: "debug"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		NewTestConnections("../../tests/contentdata.txt", STATUS, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Info(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("TransformHandler : POST should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		html, _ := ioutil.ReadFile("../../tests/htmltemplate.html")
		content, _ := ioutil.ReadFile("../../tests/contentdata.txt")
		tmp := strings.Replace(string(html), "\n", "", -1)
		h := strings.Replace(tmp, "\"", "\\\"", -1)
		//hn := strings.Replace(h, "</", "<\\/", -1)
		//hnn := strings.Replace(hn, "*/", "*// ", -1)
		hn := strings.Replace(h, "\t", " ", -1)
		//tmp = strings.Replace(string(content), "\n", "", -1)
		//c := strings.Replace(string(content), "\"", "\\\"", -1)
		//tmp := strings.Replace(c, "\n", "", -1)
		json := "{ \"file\": \"../../tests/rendered.html\",\"html\":\"" + hn + "\",\"contentdata\":" + string(content) + "}"
		logger.Info(fmt.Sprintf("json %s", json))
		req, _ := http.NewRequest("POST", "/api/v1/transform", bytes.NewBuffer([]byte(json)))
		conn := NewTestConnections("../../tests/contentdata.txt", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			TransformHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "TransformHandler", rr.Code, STATUS))
		}
	})
}
