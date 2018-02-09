package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	local_addr  = flag.String("local", "0.0.0.0:7777", "Usage: -local=<local_addr:local_port>")
	remote_addr = flag.String("remote", "", "Usage: -remote=<remote_addr:remote_port>")
)

func main() {
	var err error

	flag.Parse()

	if *local_addr == "" || *remote_addr == "" {
		flag.PrintDefaults()
		log.Fatal()
	}

	local, err := net.Listen("tcp", *local_addr)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}

	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Fatalf("accept failed: %v", err)
		}
		go forward_remote(conn, *remote_addr)
	}
}

func forward_remote(conn net.Conn, remoteAddr string) {
	defer conn.Close()

	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		return
	}

	forward(conn, remote)
}

func forward(local, remote net.Conn) {
	done := make(chan struct{}, 3)
	go func() {
		_, e := io.Copy(local, remote)
		if e != nil {
			log.Printf("remote dial: %v\n", e)
		}
		done <- struct{}{}
	}()
	go func() {
		_, e := io.Copy(remote, local)
		if e != nil {
			log.Printf("remote dial: %v\n", e)
		}
		done <- struct{}{}
	}()

	<-done

	log.Printf("remote dial end: %s %s\n", remote.RemoteAddr(), local.RemoteAddr())
}
