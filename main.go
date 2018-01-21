package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/v2af/aliyun_ddns/config"
)

var lp = log.Println

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
	alidns.NewClientWithAccessKey("", "", "")
}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
