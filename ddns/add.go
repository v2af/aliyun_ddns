package ddns

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/v2af/aliyun_ddns/config"
	"github.com/v2af/aliyun_ddns/lib"
)

func AddDomainRecord() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateAddDomainRecordRequest()
	ip, err := lib.GetIP()
	if err != nil {
		return
	}
	request.Value = ip
	request.RR = cfg.Domain.RR
	request.DomainName = cfg.Domain.DomainName
	request.Type = "A"
	rep, err := dnsClient.AddDomainRecord(request)
	if err != nil {
		log.Println("add domain record failed :", err)
		return
	}

	fmt.Println(rep)

}
