package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/v2af/file"
)

type GlobalConfig struct {
	User struct {
		REGION_ID         string `json:"region_id"`
		ACCESS_KEY_ID     string `json:"access_key_id"`
		ACCESS_KEY_SECRET string `json:"access_key_secret"`
	}
	Domain struct {
		RR         string `json:"rr"`
		DomainName string `json:"domain_name"`
		TTL        int    `json:"ttl"`
	}
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
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
