package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pingooio/stdx/httpx"
)

type Client struct {
	httpClient *http.Client
	apiToken   string
	baseURL    string
}

func NewClient(apiToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = httpx.DefaultClient()
	}

	return &Client{
		httpClient: httpClient,
		apiToken:   apiToken,
		baseURL:    "https://api.cloudflare.com",
	}
}

type requestParams struct {
	Method      string
	URL         string
	Payload     interface{}
	ServerToken *string
}

func (client *Client) request(ctx context.Context, params requestParams, dst interface{}) error {
	url := client.baseURL + params.URL
	var apiRes ApiResponse

	req, err := http.NewRequestWithContext(ctx, params.Method, url, nil)
	if err != nil {
		return err
	}

	if params.Payload != nil {
		payloadData, err := json.Marshal(params.Payload)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add(httpx.HeaderAccept, httpx.MediaTypeJson)
	req.Header.Add(httpx.HeaderContentType, httpx.MediaTypeJson)
	req.Header.Add(httpx.HeaderAuthorization, "Bearer "+client.apiToken)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		err = fmt.Errorf("cloudflare: decoding JSON body: %w", err)
		return err
	}

	if len(apiRes.Errors) != 0 {
		err = fmt.Errorf("cloudflare: %s", apiRes.Errors[0].Error())
		return err
	}

	if dst != nil {
		err = json.Unmarshal(apiRes.Result, dst)
		if err != nil {
			err = fmt.Errorf("cloudflare: decoding JSON result: %w", err)
			return err
		}
	}

	return nil
}

type ApiResponse struct {
	Result  json.RawMessage `json:"result"`
	Success bool            `json:"success"`
	Errors  []ApiError      `json:"errors"`
}

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (res ApiError) Error() string {
	return res.Message
}
