package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"yimiaoxiehou/sockstun/buffer"
)

const (
	// tcpWaitTimeout implements a TCP half-close timeout.
	tcpWaitTimeout = 60 * time.Second
)

func main() {
	proxyAddress := flag.String("proxy", "112.92.73.13:20590", "Proxy address")
	proxyUser := flag.String("user", "yimiao", "Proxy user")
	proxyPassword := flag.String("password", "c2345541z", "Proxy password")

	flag.Parse()

	dstConn := openSocks5Conn(*proxyAddress, *proxyUser, *proxyPassword)
	ifname := "MyTUN"
	ipStr := "192.168.124.1/24"

	dev, err := createTUNDevice(ifname)
	if err != nil {
		log.Fatalf("创建 TUN 设备失败: %v", err)
	}
	defer dev.Close()

	err = setIPAddress(dev, ipStr)
	if err != nil {
		log.Fatalf("设置 IP 地址失败: %v", err)
	}

	fmt.Printf("TUN 设备 '%s' 已创建，IP: %s\n", ifname, ipStr)
	pipe(&tunConn{dev}, dstConn)
}

// pipe copies data to & from provided net.Conn(s) bidirectionally.
func pipe(origin, remote net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go unidirectionalStream(remote, origin, "origin->remote", &wg)
	go unidirectionalStream(origin, remote, "remote->origin", &wg)

	wg.Wait()
}

func unidirectionalStream(dst, src net.Conn, dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := buffer.Get(buffer.RelayBufferSize)
	if _, err := io.CopyBuffer(dst, src, buf); err != nil {
		fmt.Printf("[TCP] copy data for %s: %v\n", dir, err)
	}
	buffer.Put(buf)
	// Do the upload/download side TCP half-close.
	if cr, ok := src.(interface{ CloseRead() error }); ok {
		cr.CloseRead()
	}
	if cw, ok := dst.(interface{ CloseWrite() error }); ok {
		cw.CloseWrite()
	}
	// Set TCP half-close timeout.
	dst.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
}
