package connectors

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connections struct {
	http *http.Client
	l    *simple.Logger
}

func (r *Connections) Do(req *http.Request) (*http.Response, error) {
	return r.http.Do(req)
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

// NewClientConnectors returns Connectors struct
func NewClientConnections(logger *simple.Logger) Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	conns := &Connections{l: logger, http: httpClient}
	logger.Debug(fmt.Sprintf("Connection details %v\n", conns))
	return conns
}
