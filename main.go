package main

import (
	"bytes"
	"log"
	"net"
)

var (
	maxDatagramSize            int
	multicastAddr, bouncedAddr *net.UDPAddr
)

func init() {

	maxDatagramSize = 4096

	multicastAddr = &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 251),
		Port: 5353,
	}

	bouncedAddr = &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 20), // CHANGE THIS TO THE ADDRESS YOU WANT BOUNCED
		Port: 5353,
	}
}

func main() {
	listenSocket := initListenSocket()
	dialSocket := initDialSocket()

	for {
		b := make([]byte, maxDatagramSize)
		_, srcAddr, err := listenSocket.ReadFromUDP(b)
		if err != nil {
			log.Printf("mdnsListen: ReadFrom: error %v", err)
			break
		}

		if bytes.Compare(bouncedAddr.IP, srcAddr.IP) == 0 {
			_, err := dialSocket.Write(b)
			if err != nil {
				log.Println("Write failed with:", err.Error())
			}
		}
	}

	defer listenSocket.Close()
	defer dialSocket.Close()
}

func initListenSocket() *net.UDPConn {
	socket, err := net.ListenMulticastUDP("udp4", nil, multicastAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	socket.SetReadBuffer(maxDatagramSize)
	return socket
}

func initDialSocket() *net.UDPConn {
	socket, err := net.DialUDP("udp", nil, multicastAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	return socket
}
