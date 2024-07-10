package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProxyServer struct {
	baseUrl string
	port    string
	client  *http.Client
}

func NewServer(baseUrl string, port string) *ProxyServer {
	return &ProxyServer{client: &http.Client{}, baseUrl: baseUrl, port: port}
}

func (p *ProxyServer) StartServer() {
	port := p.port
	fmt.Printf("proxy server listening on %s \n", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.requestHandler)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (p *ProxyServer) requestHandler(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()
	statusCode := http.StatusInternalServerError
	url := formFullUrl(p.baseUrl, r.URL.Path)

	proxyResp, err := proxyHandler(url, p.client, r)
	if proxyResp != nil {
		statusCode = proxyResp.StatusCode
	}

	//log response duration
	fmt.Printf("%s %s %d %v \n", r.Method, url, statusCode, time.Since(startAt))

	if proxyResp != nil {
		proxyResponseHandler(w, proxyResp)
	} else {
		internalProxyErrorHandler(w, err)
	}
}

func proxyResponseHandler(w http.ResponseWriter, proxyResp *proxyHandlerResponse) {
	setHeaders(w, proxyResp.Headers)
	w.WriteHeader(proxyResp.StatusCode)
	json.NewEncoder(w).Encode(proxyResp.Response)
}

func internalProxyErrorHandler(w http.ResponseWriter, err error) {
	jsonResponse := map[string]interface{}{
		"success": false,
		"message": err.Error(),
		"CODE":    "INTERNAL_PROXY_SERVER_ERROR",
	}

	fmt.Printf("internal proxy server error occured, please contact author. err: %s", err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(jsonResponse)
}

func setHeaders(w http.ResponseWriter, headers http.Header) {
	if headers == nil {
		return
	}

	for k, headerValues := range headers {
		for _, headerValue := range headerValues {
			w.Header().Set(k, headerValue)
		}
	}
}
