package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/aholic/ggtimer"
	"github.com/v2af/aliyun_ddns/config"
	"github.com/v2af/aliyun_ddns/ddns"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	handleConfig(*cfg)
}

func main() {
	ggtimer.NewTicker(time.Duration(5)*time.Second, func(time time.Time) {
		// ddns.ShowDomainRecordList()
		dr := ddns.GetDR()
		dr.GetRecordId()
		log.Println(dr)
		if !dr.IsExists {
			ddns.AddDomainRecord()
		} else {
			ddns.UpdateDomainRecord()
		}
	})
	select {}
	// ddns.ChangeDomainRecord()

}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
