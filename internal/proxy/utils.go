package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

func formQueryParams(query string) map[string]string {
	if query == "" {
		return nil
	}

	queryArr := strings.Split(query, "&")
	if len(queryArr) <= 0 {
		return nil
	}

	queryParamsMap := make(map[string]string, len(queryArr))
	for _, v := range queryArr {
		valueArr := strings.Split(v, "=")
		queryParamsMap[strings.TrimSpace(valueArr[0])] = strings.TrimSpace(valueArr[1])
	}

	return queryParamsMap
}

func formRequestBody(r *http.Request) (*bytes.Buffer, error) {
	body := r.Body

	if body == nil {
		return nil, nil
	}

	bodyData := make(map[string]interface{}, 0)
	err := json.NewDecoder(body).Decode(&bodyData)
	isRequestBodyEmpty := errors.Is(err, io.EOF)
	httpMethod := r.Method

	ignoredMethods := []string{http.MethodGet, http.MethodDelete}
	if err != nil && isRequestBodyEmpty && slices.Contains(ignoredMethods, httpMethod) {
		return nil, nil
	}

	if isRequestBodyEmpty && httpMethod == http.MethodPost {
		return nil, fmt.Errorf("request body is empty")
	}

	if err != nil {
		return nil, fmt.Errorf("error while decoding request body. err: %s", err.Error())
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error while parsing request body. err: %s", err.Error())
	}

	return bytes.NewBuffer(jsonData), nil
}

func formFullUrl(baseUrl string, urlPath string) string {
	return fmt.Sprintf("%s%s", strings.TrimRight(baseUrl, "/"), urlPath)
}
