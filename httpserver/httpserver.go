package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	exnet "ip"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

var (
	listenAddr string
)

func main() {

	flag.StringVar(&listenAddr, "listen-addr", ":80", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newWebserver(logger)
	go gracefullShutdown(server, logger, quit, done)

	flag.Set("v", "4")
	//glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
        }  
	<-done
	logger.Println("Server stopped")

}

func gracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func newWebserver(logger *log.Logger) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
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
	//io.WriteString(w, "======Detail of Client IP:")
	//io.WriteString(w, ip)
	//io.WriteString(w, "\n")
	fmt.Println("客户端IP：", ip)
}

