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
}

func NewServer(baseUrl string, port string) *ProxyServer {
	return &ProxyServer{baseUrl: baseUrl, port: port}
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
	client := &http.Client{}

	startAt := time.Now()
	statusCode := 200

	proxyResp, _ := proxyHandler(p.baseUrl, client, r)
	if proxyResp != nil {
		statusCode = proxyResp.StatusCode
	}

	// start log response duration
	fmt.Printf(
		"%s %s %d %v \n",
		r.Method, proxyResp.Url, statusCode, time.Since(startAt),
	)
	// end log response duration
	//fmt.Print("proxyRespErr", err)

	setHeaders(w, proxyResp.Headers)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(proxyResp.Response)
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
