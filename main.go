package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/asaskevich/EventBus"
	"github.com/v2af/aliyun_ddns/build"
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
	version := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()
	if *version {
		fmt.Println(build.String())
		os.Exit(0)
	}
	handleConfig(*cfg)
}

func main() {
	ddnsService := ddns.NewSerive()
	eventbus := EventBus.New()
	eventbus.Subscribe(lib.EventIPChange, ddnsService.OnIPChanged)
	ipservice := lib.NewIPService(eventbus)
	ipservice.Run()
}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
