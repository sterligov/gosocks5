package main

import (
	"flag"
	"github.com/sterligov/socks5"
	"strconv"
	"time"
	"os"
	"log"
)

func main() {
	logFile := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	defer f.Close()

	port := flag.String("port", "1080", "")
	username := flag.String("username", "", "")
	password := flag.String("password", "", "")
	maxConn := flag.Int("maxconn", 128, "Maximum connection number")
	timeout := flag.Int("timeout", 5000, "Timeout in milliseconds")

	flag.Parse()

	var authorizer socks5.Authorizer
	if *username != "" && *password != "" {
		authorizer = &socks5.RoughAuthorizer{
			Username: *username,
			Password: *password,
		}
	}

	s := &socks5.Socks{
		Port:           *port,
		Authorization:  authorizer,
		MaxConnections: *maxConn,
		Timeout:        *timeout,
	}
	s.Run()
}


