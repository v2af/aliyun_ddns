package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/v2af/aliyun_ddns/config"
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
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println(err)
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = cfg.Domain.DomainName
	rep, err := dnsClient.DescribeDomainRecords(request)
	if err != nil {
		log.Println(err)
	}
	for _, v := range rep.DomainRecords.Record {
		log.Println(v)
	}

}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
