package ddns

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/v2af/aliyun_ddns/config"
)

type DomainRecord struct {
	IsExists bool
	RecordId string
}

var dr *DomainRecord

func GetDR() *DomainRecord {
	if dr != nil {
		return dr
	}
	dr = &DomainRecord{}
	return dr
}

func ShowDomainRecordList() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = cfg.Domain.DomainName
	rep, err := dnsClient.DescribeDomainRecords(request)
	if err != nil {
		log.Println("get domain records failed :", err)
		return
	}
	fmt.Println("-----" + cfg.Domain.DomainName + "-----")
	for _, v := range rep.DomainRecords.Record {
		fmt.Printf("rr :%s\ntype :%s\nvalue :%s\nttl :%d\n---------------\n", v.RR, v.Type, v.Value, v.TTL)
	}
}

func (this *DomainRecord) GetRecordId() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = cfg.Domain.DomainName
	rep, err := dnsClient.DescribeDomainRecords(request)
	if err != nil {
		log.Println("get domain records failed :", err)
		return
	}
	for _, v := range rep.DomainRecords.Record {
		if cfg.Domain.RR == v.RR {
			this.IsExists = true
			this.RecordId = v.RecordId
			return
		}
	}
	this.IsExists = false
	this.RecordId = ""
}
