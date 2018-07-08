package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/asaskevich/EventBus"
	"github.com/v2af/aliyun_ddns/config"
	"github.com/v2af/aliyun_ddns/ddns"
	"github.com/v2af/aliyun_ddns/lib"
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
	ddnsService := ddns.NewDdnsSerive()
	eventbus := EventBus.New()
	eventbus.Subscribe(lib.EVENT_IP_CHANGE, ddnsService.OnIpChanged)
	ipservice := lib.NewIpService(eventbus)
	ipservice.Run()

}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
