package kit

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"gin-demo/app/manager"
)

func HttpGet(path string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	resp, err := sendHttp("GET", path, "", headers)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func HttpPost(path string, body string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	resp, err := sendHttp("POST", path, body, headers)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func HttpPatch(path string, body string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	resp, err := sendHttp("PATCH", path, body, headers)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func HttpDelete(path string, body string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	resp, err := sendHttp("DELETE", path, body, headers)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func sendHttp(method string, path string, body string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()

	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	if err != nil {
		return resp, err
	}

	setHeader(req, headers)

	manager.GetRouter().GetHandler().ServeHTTP(resp, req)
	return resp, nil
}

func setHeader(req *http.Request, headers map[string]string) {
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
