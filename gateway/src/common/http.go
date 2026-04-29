package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type APIRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
	Timeout time.Duration
}

type Response[T any] struct {
	StatusCode int
	PayLoad    *T
}

func SendRequest[T any](ctx context.Context, reqData APIRequest) (*Response[T], error) {
	var bodyReader io.Reader
	if reqData.Body != nil {
		jsonData, err := json.Marshal(reqData.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, reqData.Method, reqData.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range reqData.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: reqData.Timeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &Response[T]{
			PayLoad:    nil,
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &Response[T]{
			PayLoad:    nil,
			StatusCode: res.StatusCode,
		}, errors.New("API Request failed with status code: " + res.Status)
	}

	var result T
	if err := json.Unmarshal(resBody, &result); err != nil {
		return nil, err
	}

	return &Response[T]{
		PayLoad:    &result,
		StatusCode: res.StatusCode,
	}, nil
}
