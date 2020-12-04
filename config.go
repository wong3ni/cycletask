package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	WebConfig          `json:"web"`
	CycleTaskUnitInfos []CycleTaskUnitInfo `json:"cycletasks"`
}

var cfg *Config

func (c *Config) Load() {
	d, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("没有发现config.json文件!")
		return
	}
	err = json.Unmarshal(d, c)
	if err != nil {
		fmt.Println("解析config.json文件失败!")
		return
	}
}

func (c *Config) Save() {
	d, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println("编码成json文件出错!")
	}
	err = ioutil.WriteFile("config.json", d, 0666)
	if err != nil {
		fmt.Println("保存config.json文件出错!")
	}
}
