package main

import (
	"log"

	"cm-cloud.fr/go-pihole/db"
	"cm-cloud.fr/go-pihole/files/dnsmasq"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

var logTail *tail.Tail

// StartLogReaderService Starts service
func StartLogReaderService() {

	db.InitDataModel(dnsmasq.Log{})
	db.AutoCleanTable(dnsmasq.Log{})

	go func() {
		file := viper.GetString("dnsmasq-log-file")
		if logTail, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: true}); err == nil {
			for line := range logTail.Lines {
				db.Insert(dnsmasq.NewLog().ParseLine(line.Text))
			}
		} else {
			log.Printf("Error while tailing file %s : %s", file, err)
		}
	}()

}
