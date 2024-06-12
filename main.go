package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Create new http server that just listens on the root.
	// ensure directory exists
	err := os.Mkdir("logs", 0o755)
	if err != nil {
		fmt.Println(err)
	}

	corsHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		// Allow preflight requests to return successfully
		if req.Method == "OPTIONS" {
			return
		}

		// Call the actual handler for non-preflight requests
		handler(w, req)
	}

	http.HandleFunc("/", corsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getClientIP(r *http.Request) string {
	// Check if the IP is passed from a reverse proxy or load balancer
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Get the IP from the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error splitting host and port: %v", err)
		return r.RemoteAddr
	}
	return ip
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Dump everything about request to a file
	// create new file
	w.WriteHeader(http.StatusOK)
	ip := getClientIP(r)
	now_string := time.Now().Format("2006-01-02_15-04-05")
	err := os.Mkdir("logs/"+ip, 0o755)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create("logs/" + ip + "/" + now_string)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString("Method: " + r.Method + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("URL: " + r.URL.String() + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("Host: " + r.Host + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("RemoteAddr: " + r.RemoteAddr + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("RequestURI: " + r.RequestURI + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("Headers:\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range r.Header {
		_, _ = file.WriteString(k + " : " + v[0] + "\n")
	}
	_, err = file.WriteString("Body:\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Read body into string with ioutil.ReadAll
	body_String, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString(string(body_String))
}
