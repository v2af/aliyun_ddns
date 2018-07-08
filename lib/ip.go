package lib

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/ddliu/go-httpclient"
	"github.com/v2af/aliyun_ddns/config"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	EVENT_IP_CHANGE = "ip change"
	USERAGENT       = "aliyun_ddns"
	TIMEOUT         = 30
)

type IpService struct {
	addr       string
	port       int
	publicIP   string
	httpClient httpclient.HttpClient
	eventbus   EventBus.Bus
}

func NewIpService(eventbus EventBus.Bus) *IpService {
	s := &IpService{
		addr:     config.Config().Ip.Addr,
		port:     config.Config().Ip.Port,
		eventbus: eventbus,
	}
	s.httpClient.Defaults(httpclient.Map{
		"opt_useragent": USERAGENT,
		"opt_timeout":   TIMEOUT,
	})
	return s
}

func (this *IpService) Run() {
	ticker := time.NewTicker(time.Duration(config.Config().Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			this.handle()
		}
	}
}

func (this *IpService) handle() {
	if ip, err := this.getIP(); err != nil || ip == this.publicIP || len(ip) == 0 {
		return
	} else {
		this.eventbus.Publish(EVENT_IP_CHANGE, ip)
		this.publicIP = ip
	}

}

func (this *IpService) getIP() (ip string, err error) {
	resp, err := this.httpClient.Get(fmt.Sprintf("http://%s:%d", config.Config().Ip.Addr, config.Config().Ip.Port))
	if err != nil {
		log.Println("failed to get public ip :", err)
		return
	}
	ipstr, err := resp.ToString()
	if err != nil {
		log.Println("failed to get public ip :", err)
		return
	}
	reg := regexp.MustCompile(`(?m)^.*[\d]*\.[\d]*\.[\d]*\.[\d]*.*$`)
	str := reg.FindAllString(ipstr, -1)
	if len(str) <= 0 {
		log.Println("failed to get public ip : failed to parse response! Please check whether the URL is correct. Program is stop")
		os.Exit(0)
		return
	}
	ip = strings.TrimSpace(str[0])
	return
}
