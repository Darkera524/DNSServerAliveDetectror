package main

import (
	"github.com/toolkits/file"
	"fmt"
	"encoding/json"
	"sync"
	"time"
)

var (
	config *Config
	lock = new(sync.RWMutex)
)

type Config struct {
	Ip_list		[]string	`json:"ip_list"`
}

func ParseConfig(cfg string){
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	var c Config
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		fmt.Println(err.Error())
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
}

func GetConfig() *Config {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func CronConfig(interval int, cfg string){
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ParseConfig(cfg)
	}
}
