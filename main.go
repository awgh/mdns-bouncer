package main

import (
	"log"
	"net"

	dns "github.com/miekg/dns"
)

var (
	maxDatagramSize  int
	udpRemoteAddress *net.UDPAddr
)

func init() {

	maxDatagramSize = 4096

	udpRemoteAddress = &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 251),
		Port: 5353,
	}
}

func main() {
	listenSocket := initListenSocket()
	//dialSocket := initDialSocket()

	for {
		b := make([]byte, maxDatagramSize)
		_, cm, err := listenSocket.ReadFromUDP(b)
		if err != nil {
			log.Printf("mdnsListen: ReadFrom: error %v", err)
			break
		}
		log.Println(cm)
		msg := &dns.Msg{}
		msg.Unpack(b[:])
		log.Println(msg.String())
	}
}

func initListenSocket() *net.UDPConn {
	socket, err := net.ListenMulticastUDP("udp4", nil, udpRemoteAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	socket.SetReadBuffer(maxDatagramSize)
	return socket
}

func initDialSocket() *net.UDPConn {
	udpRemoteAddress := &net.UDPAddr{ //todo add real interface instead of nil arg
		IP:   net.IPv4(224, 0, 0, 251),
		Port: 5353,
	}
	socket, err := net.DialUDP("udp", nil, udpRemoteAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	return socket
}
