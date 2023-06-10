package udp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var clientConnection *net.UDPConn
var serverConnection *net.UDPConn

func ConnectToUDP(addr *net.UDPAddr) {
	// hier sollte etwas error handling stattfinden
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error while connecting to udp: " + err.Error())
		return
	}
	clientConnection = conn
}

func InitProxy(waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 30120,
	}
	serverConn, _ := net.ListenUDP("udp", udpAddr)
	serverConnection = serverConn

	go handleC2S()
	go handleS2C()
}

var addrToSendTo *net.UDPAddr

func handleC2S() { // FiveM -> Scriptprox -> Server
	for {
		var buf []byte = make([]byte, 1600)
		n, addr, err := serverConnection.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error while receiving C2S data: " + err.Error())
			continue
		}

		addrToSendTo = addr
		fmt.Printf("C2S msg from %s w/ size %d\n", addr, n)
		clientConnection.Write(buf[0:n])
	}
}

func handleS2C() { // Server -> Scriptprox -> FiveM
	for {
		if clientConnection == nil {
			time.Sleep(time.Second)
			continue
		}
		var buf []byte = make([]byte, 1600)
		n, _, err := clientConnection.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		fmt.Printf("S2C msg size %d\n", n)

		serverConnection.WriteTo(buf[0:n], addrToSendTo)
	}
}
