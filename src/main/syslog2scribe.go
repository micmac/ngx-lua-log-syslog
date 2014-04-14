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

type statusdesc struct {
	Blocked    int64
	Scribesent int64
}

var status statusdesc

func main() {
	config = ReadConfig("etc/s2s.ini")
	log.Println("Hello World!")
	scribequeue := make(LogStream, 100000)
	go writetoscribe(scribequeue)
	//go listenScribe("logfile.txt")

	//u = net.ListenUDP("udp4", net.UDPAddr{})
	udpaddr, err := net.ResolveUDPAddr("udp4", ":5140")
	udpsock, err := net.ListenUDP("udp4", udpaddr)
	handleError(err)
	log.Println(udpsock)
	buf := make([]byte, 4096)
	readfailures := 0
	go (func(s *statusdesc) {
		t := time.Tick(time.Second)
		for {
			<-t
			log.Printf("Sent: %d, Blocked: %d, len: %d", s.Scribesent, s.Blocked, len(scribequeue))
			s.Scribesent = 0
		}
	})(&status)

	for {
		length, err := udpsock.Read(buf)
		if err != nil {
			readfailures++
			if readfailures > 15 {
				log.Fatalf("Too many UDP read errors: %s", err)
			}
			if readfailures > 10 {
				time.Sleep(time.Second)
			}
			continue
		}
		readfailures = 0
		//log.Println(string(buf))
		select {
		case scribequeue <- LogEntry(buf[:length]):
		default:
			//log.Println("log queue blocked")
			status.Blocked++
		}
	}
}
