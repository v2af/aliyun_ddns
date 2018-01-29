package ddns

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/v2af/aliyun_ddns/config"
	"github.com/v2af/aliyun_ddns/lib"
)

func UpdateDomainRecord() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	ip, err := lib.GetIP()
	if err != nil {
		return
	}
	request.Value = ip
	request.RR = cfg.Domain.RR
	request.RecordId = GetDR().RecordId
	request.Type = "A"
	_, err = dnsClient.UpdateDomainRecord(request)
	if err != nil {
		log.Println("update domain record failed :", err)
		return
	}

}
