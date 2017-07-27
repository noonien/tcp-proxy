package main

import (
	"flag"
	"io"
	"net"
	"time"
)

var (
	listenAddr = flag.String("listen", "", "listen address")
	targetAddr = flag.String("target", "", "target address")
)

var targetIP *net.TCPAddr

func main() {
	flag.Parse()
	if len(*listenAddr) == 0 || len(*targetAddr) == 0 {
		flag.Usage()
		return
	}

	resolveTargetIP()
	go func() {
		ticker := time.NewTicker(10 * time.Minute)

		for range ticker.C {
			resolveTargetIP()
		}
	}()

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConn(conn)
	}
}

func resolveTargetIP() {
	var err error
	targetIP, err = net.ResolveTCPAddr("tcp", *targetAddr)
	if err != nil {
		panic(err)
	}
}

func handleConn(src net.Conn) {
	defer src.Close()

	dst, err := net.DialTCP("tcp", nil, targetIP)
	if err != nil {
		return
	}
	defer dst.Close()

	errc := make(chan error, 1)
	go proxy(errc, src, dst)
	go proxy(errc, dst, src)
	<-errc
}

func proxy(errc chan error, dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	errc <- err
}
