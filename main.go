package main

import (
	"flag"
	"log"
	"runtime"

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
	ddns.ShowDomainRecordList()

}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
