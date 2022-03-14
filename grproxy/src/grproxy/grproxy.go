package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/go-zookeeper/zk"
	"time"
)

var serverCount = 0

// These constant is used to define server
const (
	SERVER1    = "http://nginx:80"   	// nginx
	SERVER2    = "http://gserve1:7000" 	// gserve1
	SERVER3    = "http://gserve2:7000"    // gserve2
	PROXY_PORT = "3000"
)

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Log the typeform payload and redirect url
func logRequestPayload(proxyURL string) {
	log.Printf("proxy_url: %s\n", proxyURL)
}

// forward to servers
func getProxyURL() string {
	// round robin
	conn, _, _ := zk.Connect([]string{"zookeeper:2181"}, time.Second)
	var connected_servers []string 

	exists1, _, _ := conn.Exists("/gserve1")
	if exists1 == true {
		connected_servers = append(connected_servers, SERVER2)
	}

	exists2, _, _ := conn.Exists("/gserve2")
	if exists2 == true {
		connected_servers = append(connected_servers, SERVER3)
	}

    // if both gserve is down, do something
	if len(connected_servers) == 0 {
		return SERVER1
	}

	server := connected_servers[serverCount]
	serverCount++

	// reset the counter and start from the beginning
	if serverCount >= len(connected_servers) {
		serverCount = 0
	}

	return server
}

// Given a request send it to the appropriate url
func handleRequestForNginx(res http.ResponseWriter, req *http.Request) {
	// get url
	log.Printf("nginx\n")
	url := SERVER1

	// Just logging
	logRequestPayload(url)

	// Actual work being done here
	serveReverseProxy(url, res, req)
}

func handleRequestForGserve(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "POST" {
		log.Printf("Method is not GET or POST. It's %s\n", req.Method)
		return
	}
	// get url
	url := getProxyURL()

	// Just logging
	logRequestPayload(url)

	// Actual work being done here
	serveReverseProxy(url, res, req)
}

func main() {
	// start server
	http.HandleFunc("/", handleRequestForNginx)
	http.HandleFunc("/library", handleRequestForGserve)

	if err := http.ListenAndServe(":"+PROXY_PORT, nil); err != nil {
    	log.Fatal("ListenAndServe: ", err)
	}
}

// https://medium.com/swlh/proxy-server-in-golang-43e2365d9cbc
