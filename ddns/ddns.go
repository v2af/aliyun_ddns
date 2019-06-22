package ddns

import (
	"fmt"
	"log"
	"time"

	"github.com/Rican7/retry/strategy"

	"github.com/Rican7/retry"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/pkg/errors"
	"github.com/v2af/aliyun_ddns/config"
)

var (
	// ErrorDomainRecordIsExists ErrorDomainRecordIsExists
	ErrorDomainRecordIsExists = errors.New("domain record is exists")
	// ErrorDomainRecordIsNotUpdate ErrorDomainRecordIsNotUpdate
	ErrorDomainRecordIsNotUpdate = errors.New("domain record is not update")
)

// Service Service
type Service struct {
	recordID string
	publicIP string
}

// NewSerive NewSerive
func NewSerive() *Service {
	return &Service{}
}

// OnIPChanged OnIPChanged
func (s *Service) OnIPChanged(ip string) {
	s.publicIP = ip
	switch err := s.showCurrentDomainRecord(); err {
	case ErrorDomainRecordIsExists:
		if err := retry.Retry(
			func(attempt uint) error {
				if err := s.UpdateDomainRecord(); err != nil {
					return err
				}
				return nil
			},
			strategy.Limit(5),
			strategy.Delay(3*time.Second)); err != nil {
			// smtp
			log.Println(err)
		}
		break
	case nil:
		if err := retry.Retry(
			func(attempt uint) error {
				if err := s.AddDomainRecord(); err != nil {
					return err
				}
				return nil
			},
			strategy.Limit(5),
			strategy.Delay(3*time.Second)); err != nil {
			// smtp
			log.Println(err)
		}
		break
	default:
		log.Println(err)
	}
}

func (s *Service) showCurrentDomainRecord() error {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.RegionID, cfg.User.AccessKeyID, cfg.User.AccessKeySecret)
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
			s.recordID = v.RecordId
			if s.publicIP == v.Value {
				return ErrorDomainRecordIsNotUpdate
			}
			return ErrorDomainRecordIsExists
		}
	}
	return nil
}

// AddDomainRecord AddDomainRecord
func (s *Service) AddDomainRecord() error {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.RegionID, cfg.User.AccessKeyID, cfg.User.AccessKeySecret)
	if err != nil {
		log.Println("create dns client failed :", err)
		return err
	}
	request := alidns.CreateAddDomainRecordRequest()
	request.Value = s.publicIP
	request.RR = cfg.Domain.RR
	request.DomainName = cfg.Domain.DomainName
	request.Type = "A"
	request.TTL = requests.NewInteger(cfg.Domain.TTL)
	_, err = dnsClient.AddDomainRecord(request)
	if err != nil {
		log.Println("add domain record failed :", err)
		return err
	}
	return nil
}

// UpdateDomainRecord UpdateDomainRecord
func (s *Service) UpdateDomainRecord() error {
	cfg := config.Config()
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.User.RegionID, cfg.User.AccessKeyID, cfg.User.AccessKeySecret)
	if err != nil {
		log.Println("create dns client failed :", err)
		return err
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Value = s.publicIP
	request.RR = cfg.Domain.RR
	request.RecordId = s.recordID
	request.Type = "A"
	request.TTL = requests.NewInteger(cfg.Domain.TTL)

	_, err = dnsClient.UpdateDomainRecord(request)
	if err != nil {
		log.Println("update domain record failed :", err)
		return err
	}
	return nil
}
