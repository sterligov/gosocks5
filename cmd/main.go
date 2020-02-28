package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sterligov/socks5"
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
	timeout := flag.Int("deadline", 120, "Read-write deadline in seconds")

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
		Deadline:       time.Duration(*timeout) * time.Second,
	}
	s.Run()
}
