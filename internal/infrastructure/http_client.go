package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type httpClient struct {
}

type HttpClient interface {
	GetRequest(ctx context.Context, url string, authorization string) ([]byte, error)
	PostRequest(ctx context.Context, url string, reqBody any, authorization string) ([]byte, error)
	UploadFile(ctx context.Context, endpoint, filepathName string, authorization string) ([]byte, error)
}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

func (c *httpClient) PostRequest(ctx context.Context, url string, reqBody any, authorization string) ([]byte, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(jsonBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("%s", string(bodyBytes))
	}
	return bodyBytes, nil
}

func (c *httpClient) GetRequest(ctx context.Context, url string, authorization string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		errResp := ErrResp{}
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("code=%s: message=%s", errResp.Code, errResp.Message)
	}
	return bodyBytes, nil
}

func (c *httpClient) UploadFile(ctx context.Context, endpoint, filepathName string, authorization string) ([]byte, error) {
	file, err := os.Open(filepathName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filepathName))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		errResp := ErrResp{}
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("code=%s: message=%s", errResp.Code, errResp.Message)
	}
	return bodyBytes, nil
}

type ErrResp struct {
	Code    string
	Message string
}
