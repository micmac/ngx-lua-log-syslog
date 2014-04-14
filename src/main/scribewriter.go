// scribewriter
package main

import (
	"fmt"
	"github.com/prezi/go-thrift/examples/scribe"
	"github.com/prezi/go-thrift/thrift"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type ScribeLogger struct {
	ScribeClient scribe.ScribeClient
	Category     string
}

func (S ScribeLogger) FormatLog(logentry LogEntry) (logline string) {
	//scribe.Log(fmt.Sprintf("now hostname metricstream hostid:checkname:tstamp:values:msg%s\n", firstline))
	const timeformat = "2006-01-02 15:04:05"

	now := time.Now()
	myname, _ := os.Hostname()
	logline = fmt.Sprintf("%s,%03d %s %s %s",
		now.Format(timeformat), now.Nanosecond()/1e6,
		myname,
		S.Category,
		logentry,
	)
	return

}

func (S ScribeLogger) Log(msg string) error {
	var logentry = scribe.LogEntry{S.Category, strings.TrimRight(msg, "\n") + "\n"}
	_, err := S.ScribeClient.Log([]*scribe.LogEntry{&logentry})
	status.Scribesent++
	return err
}

func NewScribeLogger(scribeserver string, category string) (logger *ScribeLogger) {
	logger = new(ScribeLogger)
	conn, err := net.Dial("tcp", scribeserver)
	if err != nil {
		log.Printf("Can't connect to scribe server %s: %s", scribeserver, err)
		return nil
	}

	logger.ScribeClient.Client = thrift.NewClient(
		thrift.NewFramedReadWriteCloser(conn, 0),
		thrift.NewBinaryProtocol(true, false),
		false)
	logger.Category = category
	log.Printf("Connected to scribe server %s, category %s: %s,%s", scribeserver, logger.Category, conn, err)
	return
}
