package ddns

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/pkg/errors"
	"github.com/v2af/aliyun_ddns/config"
	"log"
)

var ERROR_DOMAIN_RECORD_IS_EXISTS = errors.New("domain record is exists")
var ERROR_DOMAIN_RECORD_IS_NOT_UPDATE = errors.New("domain record is not update")

type DdnsService struct {
	recordId string
	publicIP string
}

func NewDdnsSerive() *DdnsService {
	return &DdnsService{}
}

func (this *DdnsService) OnIpChanged(ip string) {
	this.publicIP = ip
	switch err := this.showCurrentDomainRecord(); err {
	case ERROR_DOMAIN_RECORD_IS_EXISTS:
		this.UpdateDomainRecord()
		break
	case nil:
		this.AddDomainRecord()
		break
	}
}

func (this *DdnsService) showCurrentDomainRecord() error {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return err
	}
	request := alidns.CreateDescribeDomainRecordsRequest()

	request.DomainName = cfg.Domain.DomainName
	rep, err := dnsClient.DescribeDomainRecords(request)
	if err != nil {
		log.Println("get domain records failed :", err)
		return err
	}
	for _, v := range rep.DomainRecords.Record {
		if v.RR == cfg.Domain.RR {
			fmt.Printf("current domain record :\nrr :%s\ntype :%s\nvalue :%s\nttl :%d\n---------------\n", v.RR, v.Type, v.Value, v.TTL)
			this.recordId = v.RecordId
			if this.publicIP == v.Value {
				return ERROR_DOMAIN_RECORD_IS_NOT_UPDATE
			}
			return ERROR_DOMAIN_RECORD_IS_EXISTS
		}
	}
	return nil
}

func (this *DdnsService) AddDomainRecord() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateAddDomainRecordRequest()

	request.Value = this.publicIP
	request.RR = cfg.Domain.RR
	request.DomainName = cfg.Domain.DomainName
	request.Type = "A"
	_, err = dnsClient.AddDomainRecord(request)
	if err != nil {
		log.Println("add domain record failed :", err)
		return
	}
}

func (this *DdnsService) UpdateDomainRecord() {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.REGION_ID, cfg.User.ACCESS_KEY_ID, cfg.User.ACCESS_KEY_SECRET)
	if err != nil {
		log.Println("create dns client failed :", err)
		return
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	if err != nil {
		return
	}
	request.Value = this.publicIP
	request.RR = cfg.Domain.RR
	request.RecordId = this.recordId
	request.Type = "A"
	_, err = dnsClient.UpdateDomainRecord(request)
	if err != nil {
		log.Println("update domain record failed :", err)
		return
	}

}
