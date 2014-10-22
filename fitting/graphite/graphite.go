package graphite

import (
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"net"
	"time"
	"strconv"
	"strings"
)

type GraphiteFitting struct {
	port     uint
	writer   engine.Writer
	listener net.Listener
}

func (self *GraphiteFitting) Name() string {
	return "Graphite"
}

func (self *GraphiteFitting) Start(app app.App) {
	self.writer = app.GetWriter()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", self.port))
	if err != nil {
		log.Error("Error start Graphite fitting %s", err.Error())
		return
	}
	self.listener = listener
	for {
		conn, err := self.listener.Accept()
		if err != nil {
			return
		}
		go self.handleConnection(conn)
	}

}

func (self *GraphiteFitting) Stop() {
	self.listener.Close()
}

func (self *GraphiteFitting) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			log.Warn("Malformed graphite metric data '%s'", line)
			continue
		}
		metricID := metric.MetricID(parts[0])
		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Warn("Error parsing value as float %s", parts[1])
			continue
		}
		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Warn("Error parsing timestamp %s", parts[2])
			continue
		}
		tm := time.Unix(timestamp, 0).UTC()

		self.writer.Insert(tm, metricID, value)
	}
	if err := scanner.Err(); err != nil {
		log.Warn("Error reading data from conn %v: %s:", conn, err.Error())
	}
	conn.Close()
}
