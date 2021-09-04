package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type ProveResponse struct {
	Servertime   int64  `json:"serverTime"`
	Salt         string `json:"salt"`
	Accesstoken  string `json:"accessToken"`
	Tokentype    string `json:"tokenType"`
	Refreshtoken string `json:"refreshToken"`
	Salt1        int    `json:"salt1"`
	Salt2        int    `json:"salt2"`
	Salt3        int    `json:"salt3"`
	Salt4        int    `json:"salt4"`
	Salt5        int    `json:"salt5"`
}

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"
)

type Client struct {
	httpClient *http.Client
	BaseURL    *url.URL
	UserAgent  string
	Headers    string
}

func NewClient(httpClient *http.Client, apiURL string, auth string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(apiURL)

	c := &Client{
		httpClient: httpClient,
		UserAgent:  defaultUserAgent,
		BaseURL:    baseURL,
		Headers:    auth,
	}
	return c
}

// Do sends an API request and returns the API response. The API response is JSON
// decoded and stored in the value pointed to by v.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
		return nil, err
	}
	return resp, nil
}

// NewRequest creates an API request. The given URL is relative to the Client's
// BaseURL.
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {

	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	if c.Headers != "" {
		req.Header.Set("Authorization", c.Headers)
	}
	return req, nil
}

func checkResponse(r *http.Response) error {
	status := r.StatusCode
	if status >= 200 && status <= 299 {
		return nil
	}

	return fmt.Errorf("request failed with status %d", status)
}

func (c *Client) Wasm(prove ProveResponse) ProveResponse {
	wasmBytes, _ := ioutil.ReadFile("css.wasm")

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, _ := wasmer.NewModule(store, wasmBytes)

	importObject := wasmer.NewImportObject()
	instance, _ := wasmer.NewInstance(module, importObject)
	cdx, _ := instance.Exports.GetFunction("cdx")
	rdx, _ := instance.Exports.GetFunction("rdx")

	n, _ := cdx(prove.Salt1, prove.Salt2, prove.Salt3, prove.Salt4, prove.Salt5)
	l, _ := rdx(prove.Salt1, prove.Salt2, prove.Salt4, prove.Salt3, prove.Salt5)
	i, _ := cdx(prove.Salt2, prove.Salt1, prove.Salt3, prove.Salt5, prove.Salt4)
	r, _ := rdx(prove.Salt2, prove.Salt1, prove.Salt3, prove.Salt4, prove.Salt5)

	prove.Accesstoken = prove.Accesstoken[:n.(int32)] + prove.Accesstoken[n.(int32)+1:l.(int32)] + prove.Accesstoken[l.(int32)+1:]
	prove.Refreshtoken = prove.Refreshtoken[:i.(int32)] + prove.Refreshtoken[i.(int32)+1:r.(int32)] + prove.Refreshtoken[r.(int32)+1:]

	prove.Accesstoken = fmt.Sprintf("Salter %v", prove.Accesstoken)
	return prove
}
