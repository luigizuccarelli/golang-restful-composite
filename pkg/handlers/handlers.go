package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-restful-composite/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-restful-composite/pkg/schema"
)

const (
	CONTENTTYPE      string = "Content-Type"
	APPLICATIONJSON  string = "application/json"
	EXECUTECOMPOSITE string = "ExecuteComposite"
	EXECUTEREQUEST   string = "ExecuteRequest"
)

func ExecuteRequest(data schema.Composite, ch chan (schema.Composite), con connectors.Clients) {
	var req *http.Request
	var err error

	data.Status = "OK"
	data.StatusCode = "200"
	data.Message = "Request called successfully"
	if data.Method != "POST" {
		req, err = http.NewRequest(data.Method, data.Url, nil)
	} else {
		req, err = http.NewRequest(data.Method, data.Url, strings.NewReader(data.Payload))
	}
	if err != nil {
		con.Error(fmt.Sprintf(EXECUTEREQUEST+" %v", err))
		data.Status = "KO"
		data.StatusCode = "500"
		data.Message = fmt.Sprintf(" %v", err)
		ch <- data
		return
	}
	// add headers if set
	if data.Headers != "" {
		headers := strings.Split(data.Headers, ",")
		for i, _ := range headers {
			kv := strings.Split(headers[i], ":")
			req.Header.Set(kv[0], kv[1])
		}
	}
	resp, err := con.Do(req)
	con.Info(fmt.Sprintf(EXECUTEREQUEST+" connected to host %s", data.Url))
	if err != nil || resp.StatusCode != 200 {
		con.Error(fmt.Sprintf(EXECUTEREQUEST+" %v", err))
		data.Status = "KO"
		data.StatusCode = "500"
		data.Message = fmt.Sprintf(" %v", err)
		ch <- data
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		con.Error(fmt.Sprintf(EXECUTEREQUEST+" %v", err))
		data.Status = "KO"
		data.StatusCode = "500"
		data.Message = fmt.Sprintf(" %v", err)
		ch <- data
		return
	}
	data.Payload = string(body)
	ch <- data
	con.Debug(fmt.Sprintf(EXECUTEREQUEST+" channel data %v", ch))
}

func ExecuteComposite(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var response *schema.Response
	var dataIn schema.RequestSchema
	var responseData schema.Composite
	var payload schema.SchemaInterface
	var bError bool = false

	addHeaders(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = &schema.Response{StatusCode: "500", Status: "KO", Message: fmt.Sprintf(EXECUTECOMPOSITE+" could not read body data %v\n", err), Results: payload}
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		e := json.Unmarshal(body, &dataIn)
		if e != nil {
			response = &schema.Response{StatusCode: "500", Status: "KO", Message: fmt.Sprintf(EXECUTECOMPOSITE+" could not unmarshal data %v\n", e), Results: payload}
			w.WriteHeader(http.StatusInternalServerError)
			con.Error(EXECUTECOMPOSITE+" %v", e)
		} else {
			// iterate through each object
			var ch = make(chan schema.Composite)
			for x, _ := range dataIn.Requests {
				go ExecuteRequest(dataIn.Requests[x], ch, con)
				responseData = <-ch
				payload.Requests = append(payload.Requests, responseData)
				payload.MergedContent = append(payload.MergedContent, responseData.Payload)
				if dataIn.Strategy == "failonce" && responseData.Status == "KO" {
					bError = true
					break
				}
			}
			close(ch)
			payload.MetaInfo = dataIn.MetaInfo
			payload.Strategy = dataIn.Strategy
			payload.LastUpdate = time.Now().UnixNano()
			if bError {
				response = &schema.Response{StatusCode: "500", Status: "KO", Message: EXECUTECOMPOSITE + " failed on request", Results: payload}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				response = &schema.Response{StatusCode: "200", Status: "OK", Message: EXECUTECOMPOSITE + " data successfully aggregated", Results: payload}
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
