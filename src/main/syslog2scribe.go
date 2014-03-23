// ms project main.go
package main

import (
	"log"
	"net"
	"time"
)

type LogEntry string
type LogStream chan LogEntry

func handleError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func writetoscribe(queue LogStream) {
	const MaxRetries = 3
	const RetrySleep = 10
	var S *ScribeLogger = nil
	retries := MaxRetries
	for {
		for S == nil {
			retries--
			S = NewScribeLogger("localhost:1463", "syslog")
			if S != nil {
				break
			}
			if retries == 0 {
				time.Sleep(RetrySleep * time.Second)
				retries = MaxRetries
			} else {
				time.Sleep(time.Second)
			}
		}
		logline := <-queue
		err := S.Log(S.FormatLog(logline))
		if err != nil {
			log.Printf("WtS: %s", err)
			S = nil
		}
	}
}

func main() {
	config = ReadConfig("etc/s2s.ini")
	log.Println("Hello World!")
	scribequeue := make(LogStream)
	go writetoscribe(scribequeue)

	//u = net.ListenUDP("udp4", net.UDPAddr{})
	udpaddr, err := net.ResolveUDPAddr("udp4", ":5140")
	udpsock, err := net.ListenUDP("udp4", udpaddr)
	handleError(err)
	log.Println(udpsock)
	buf := make([]byte, 4096)
	for {
		udpsock.Read(buf)
		log.Println(string(buf))
		scribequeue <- LogEntry(buf)
	}

}
