package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Requests struct {
	headers           map[string]string
	extra             string
	arbitrary_allowed bool
}

func NewRequests(headers map[string]string) *Requests {
	return &Requests{
		headers:           headers,
		extra:             "forbid",
		arbitrary_allowed: true,
	}
}

func (r *Requests) applyHeaders(req *http.Request) {
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
}

func (r *Requests) get(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	r.applyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (r *Requests) post(url string, data map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	r.applyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (r *Requests) patch(url string, data map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	r.applyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (r *Requests) put(url string, data map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	r.applyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (r *Requests) delete(url string) (string, error) {
	req, _ := http.NewRequest("DELETE", url, nil)
	r.applyHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
