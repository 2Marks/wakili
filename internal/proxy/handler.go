package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/2marks/wakili/internal/errs"
)

type proxyHandlerResponse struct {
	Url        string
	Response   map[string]interface{}
	StatusCode int
	Headers    http.Header
}

func proxyHandler(url string, client *http.Client, r *http.Request) (*proxyHandlerResponse, error) {
	//handle request body
	requestBody, err := formRequestBody(r)
	if err != nil {
		return nil, errs.ErrInternal(err)
	}

	req, err := getNewRequest(r.Method, url, requestBody)
	if err != nil {
		return nil, errs.ErrInternal(err)
	}

	/**** start handle query params *****/
	queryParams := formQueryParams(r.URL.RawQuery)
	if queryParams != nil {
		q := req.URL.Query()

		for k, v := range queryParams {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}
	/**** end handle query params *****/

	/*** start request headers **/
	if r.Header != nil {
		for k, v := range r.Header {
			for _, headerVal := range v {
				req.Header.Add(k, headerVal)
			}
		}
	}
	/** end request headers **/

	resp, err := client.Do(req)
	if err != nil {
		return nil, errs.ErrInternal(err)
	}

	proxyHandlerResp := proxyHandlerResponse{
		Url:        url,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	defer resp.Body.Close()

	httpResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &proxyHandlerResp, err
	}

	if string(httpResponseBody) != "" {
		err = json.Unmarshal(httpResponseBody, &proxyHandlerResp.Response)
		if err != nil {
			return &proxyHandlerResp, err
		}
	}

	return &proxyHandlerResp, nil
}

func getNewRequest(method string, url string, requestBody *bytes.Buffer) (*http.Request, error) {
	if requestBody != nil {
		return http.NewRequest(method, url, requestBody)
	}

	return http.NewRequest(method, url, nil)
}
