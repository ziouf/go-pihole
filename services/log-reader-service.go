package services

import (
	"log"

	"cm-cloud.fr/go-pihole/models"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

// StartLogReaderService Starts service
func StartLogReaderService() {

	go func() {
		file := viper.GetString("dnsmasq-log-file")
		if t, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: true}); err == nil {
			for l := range t.Lines {
				models.Insert(models.NewDnsmasqLog().ParseLine(l.Text))
			}
		} else {
			log.Printf("Error while tailing file %s : %s", t.Filename, err)
		}
	}()

}
