package main

import (
	"flag"
	"fmt"
	"io"
	exnet "ip"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	flag.Set("v", "4")
	//glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	//写入服务端的状态码随便写了一个状态码
	//w.WriteHeader(http.StatusContinue)
	//fmt.Fprintln(w, "状态码正常，能够访问")
	//返回客户端的头
	io.WriteString(w, "======Detail of the request header :=====\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	//返回version 环境变量
	io.WriteString(w, "======Detail of Os VERSION:\n")
	io.WriteString(w, os.Getenv("VERSION"))
	io.WriteString(w, "\n")

	//返回客户端IP
	ip := exnet.ClientPublicIP(r)
	if ip == "" {
		ip = exnet.ClientIP(r)
	}
	io.WriteString(w, "======Detail of Client IP:")
	io.WriteString(w, ip)
	io.WriteString(w, "\n")
}

