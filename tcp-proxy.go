package main

import (
	"flag"
	"io"
	"net"
)

var (
	listenAddr = flag.String("listen", "", "listen address")
	targetAddr = flag.String("target", "", "target address")
)

func main() {
	flag.Parse()
	if len(*listenAddr) == 0 || len(*targetAddr) == 0 {
		flag.Usage()
		return
	}

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

func handleConn(src net.Conn) {
	defer src.Close()

	dst, err := net.Dial("tcp", *targetAddr)
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
