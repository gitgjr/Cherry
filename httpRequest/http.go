package httpRequest

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendPostRequest(data any, url string) (*http.Response, error) {
	bodyReader, err := anyToJsonToReader(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func SendPutRequest(data any, url string) (*http.Response, error) {
	bodyReader, err := anyToJsonToReader(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func SendFileRequest(data any, url string) (*http.Response, error) {

	bodyReader, err := anyToJsonToReader(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func anyToJsonToReader(data any) (*bytes.Reader, error) {
	JsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(JsonData)
	return bodyReader, nil
}
