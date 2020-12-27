package httpRequest

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpHeaders map[string]string

func PostData(url string, body []byte, headers HttpHeaders, timeoutInSec int) (code int, response []byte, err error) {
	clientHTTP := http.Client{}

	clientHTTP.Timeout = time.Second * 60 //set timeout to 1 minute (default)
	if timeoutInSec != 0 {
		clientHTTP.Timeout = time.Second * time.Duration(timeoutInSec) //set timeout to custom timeout
	}

	req, err := http.NewRequest("POST",
		url, bytes.NewReader(body))

	if err != nil {
		return 0, nil, err
	}

	//add header
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	//execute http post
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return 0, nil, nil
	}
	defer resp.Body.Close()

	//read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, responseBody, nil
}

func PostDataWithBasicAuth(url string, username, password string, body []byte, headers HttpHeaders, timeoutInSec int) (code int, response []byte, err error) {
	//build basic auth token
	authToken := base64.URLEncoding.EncodeToString([]byte(username + ":" + password))
	if headers == nil {
		headers = make(map[string]string)
	}

	//add Authorization header
	headers["Authorization"] = "Basic " + authToken

	//check if Client-Id header is not exists
	//this header is used in some authorization server
	if _, ok := headers["Client-Id"]; !ok {
		headers["Client-Id"] = username
	}

	clientHTTP := http.Client{}

	clientHTTP.Timeout = time.Second * 60 //set timeout to 1 minute (default)
	if timeoutInSec != 0 {
		clientHTTP.Timeout = time.Second * time.Duration(timeoutInSec) //set timeout to custom timeout
	}

	req, err := http.NewRequest("POST",
		url, bytes.NewReader(body))

	if err != nil {
		return 0, nil, err
	}

	//add header
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	//execute http post
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return 0, nil, nil
	}
	defer resp.Body.Close()

	//read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, responseBody, nil
}
