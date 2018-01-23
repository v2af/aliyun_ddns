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
		return fmt.Errorf("使用 -c 指定配置文件")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("配置文件%s不存在", cfg)
	}

	ConfigFile = cfg
	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return fmt.Errorf("读取配置文件 %s 失败,原因:  %s", cfg, err.Error())
	}
	var c GlobalConfig
	err = json.Unmarshal(data, &c)
	if err != nil {
		return fmt.Errorf("解析配置文件 %s 失败,原因: %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("读取配置文件", cfg, "成功")
	return nil
}
