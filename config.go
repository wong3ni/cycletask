package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	WebConfig          `json:"web"`
	LoggerInfo         `json:"log"`
	CycleTaskUnitInfos []CycleTaskUnitInfo `json:"cycletasks"`
}

var cfg *Config

func (c *Config) Load() {
	abspath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	d, err := ioutil.ReadFile(filepath.Join(abspath, "config.json"))
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
	abspath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err = ioutil.WriteFile(filepath.Join(abspath, "config.json"), d, 0666)
	if err != nil {
		fmt.Println("保存config.json文件出错!")
	}
}
