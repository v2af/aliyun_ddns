package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/v2af/file"
)

// GlobalConfig GlobalConfig
type GlobalConfig struct {
	Interval int `json:"interval" default:"5" envconfig:"INTERVAL"`
	User     struct {
		RegionID        string `json:"region_id" default:"cn-hangzhou" envconfig:"ALIYUN_REGION_ID"`
		AccessKeyID     string `json:"access_key_id" required:"true" envconfig:"ALIYUN_ACCESS_KEY_ID"`
		AccessKeySecret string `json:"access_key_secret" required:"true" envconfig:"ALIYUN_ACCESS_KEY_SECRET"`
	} `json:"user"`
	Domain struct {
		RR         string `json:"rr" default:"ddns" envconfig:"DOMAIN_RR"`
		DomainName string `json:"domain_name" required:"true" envconfig:"DOMAIN_NAME"`
		TTL        int    `json:"ttl" default:"600" envconfig:"DNS_TTL"`
	} `json:"domain"`
	IP struct {
		Addr string `json:"addr" default:"script.v2af.com" envconfig:"GET_PUBLIC_IP_ADDR"`
		Port int    `json:"port" default:"3000" envconfig:"GET_PUBLIC_IP_PORT"`
	} `json:"ip"`
}

var (
	// ConfigFile ConfigFile
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config Config
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// Parse Parse
func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return envConfig()
	}

	ConfigFile = cfg
	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}
	var c GlobalConfig
	err = json.Unmarshal(data, &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}

func envConfig() error {
	var c GlobalConfig
	err := envconfig.Process("ddns", &c)
	if err != nil {
		log.Println(err)
		return err
	}
	config = &c
	log.Println("load environment successfully")
	return nil
}
