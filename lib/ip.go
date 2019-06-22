package lib

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/ddliu/go-httpclient"
	"github.com/v2af/aliyun_ddns/config"
)

const (
	// EventIPChange EventIPChange
	EventIPChange = "ip change"
	// USERAGENT USERAGENT
	USERAGENT = "aliyun_ddns"
	// TIMEOUT TIMEOUT
	TIMEOUT = 30
)

// IPService IPService
type IPService struct {
	addr       string
	port       int
	publicIP   string
	httpClient httpclient.HttpClient
	eventbus   EventBus.Bus
}

// NewIPService NewIPService
func NewIPService(eventbus EventBus.Bus) *IPService {
	s := &IPService{
		addr:     config.Config().IP.Addr,
		port:     config.Config().IP.Port,
		eventbus: eventbus,
	}
	s.httpClient.Defaults(httpclient.Map{
		"opt_useragent": USERAGENT,
		"opt_timeout":   TIMEOUT,
	})
	return s
}

// Run Run
func (s *IPService) Run() {
	ticker := time.NewTicker(time.Duration(config.Config().Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			s.handle()
		}
	}
}

func (s *IPService) handle() {
	if ip, err := s.getIP(); err != nil || ip == s.publicIP || len(ip) == 0 {
		return
	} else {
		s.eventbus.Publish(EventIPChange, ip)
		s.publicIP = ip
	}

}

func (s *IPService) getIP() (ip string, err error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("http://%s:%d", config.Config().IP.Addr, config.Config().IP.Port))
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
