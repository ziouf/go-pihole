package services

import (
	"log"
	"runtime"
	"sync"

	"cm-cloud.fr/go-pihole/models"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

var logChan = make(chan *models.DnsmasqLog)

func StartLogReaderService() {
	file := viper.GetString("dnsmasq-log-file")
	if t, err := tail.TailFile(file, tail.Config{ /* Follow: true, ReOpen: true */ }); err == nil {

		for i := 0; i < runtime.NumCPU()*4; i++ {
			go logFileReader(t)
		}
		for i := 0; i < 4; i++ {
			go publishLogs()
		}

	} else {
		log.Printf("Error while tailing file %s : %s", t.Filename, err)
	}
}

func logFileReader(t *tail.Tail) {
	for line := range t.Lines {
		// log.Println("Reading line", line)
		l := &models.DnsmasqLog{}
		l.ParseLine(line.Text)
		logChan <- l
	}
}

const bulkSize = 250

var mutex sync.Mutex

func publishLogs() {
	buffer := make([]*models.DnsmasqLog, 0)
	for l := range logChan {
		buffer = append(buffer, l)

		if len(buffer) >= bulkSize {
			buffer = persist(buffer)
		}
	}
}

func persist(buffer []*models.DnsmasqLog) []*models.DnsmasqLog {
	mutex.Lock()
	tx := models.Db.Begin()
	for _, item := range buffer {
		if tx.NewRecord(item) {
			tx.Create(item)
		}
	}
	tx.Commit()
	mutex.Unlock()
	return make([]*models.DnsmasqLog, 0)
}
