package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"
)

func main() {
	addr := flag.String("addr", "localhost:8000", "listen address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/t", handle)
	http.ListenAndServe(*addr, mux)
}

type echoJSON struct {
	Echo       string              `json:"echo"`
	RemoteAddr string              `json:"remote_addr"`
	Header     map[string][]string `json:"header"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	const XPetRemoteAddr = "X-Pet-Remote-Addr"
	var rspJSON echoJSON
	rspJSON.Echo = r.URL.Query().Get("echo")
	rspJSON.RemoteAddr = r.Header.Get(XPetRemoteAddr)
	if len(rspJSON.RemoteAddr) == 0 {
		rspJSON.RemoteAddr = r.RemoteAddr[:strings.IndexByte(r.RemoteAddr, ':')]
	}
	rspJSON.Header = make(map[string][]string)
	for k, v := range r.Header {
		if k == XPetRemoteAddr {
			continue
		}
		rspJSON.Header[k] = v
	}
	rspData, _ := json.Marshal(rspJSON)
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	w.Write(rspData)
}
