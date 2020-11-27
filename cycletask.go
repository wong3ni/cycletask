package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

var cycletaskmap map[string]*CycleTask

func HandleInput(ctx context.Context, cancel context.CancelFunc) {
	var cmd string
	var tag string
	var req_url string
	var des_url string
	des_url = "http://192.168.168.103:8080/boat"
	for {
		fmt.Print(">")
		fmt.Scanln(&cmd, &tag, &req_url, &des_url)
		switch cmd {
		case "add":
			cycletask := NewCycleTask(req_url, des_url, tag, 2)
			cycletaskmap[tag] = cycletask
			fmt.Println("OK!")
		case "start":
			cycletask := cycletaskmap[tag]
			cycletask.StartCycle(ctx)
			fmt.Println("OK!")
		case "del":
			cycletask := cycletaskmap[tag]
			cycletask.StopCycle()
			cycletask.Wait()
			delete(cycletaskmap, tag)
			fmt.Println("OK!")
		case "view":
			for k, v := range cycletaskmap {
				fmt.Println("tag:", k, " req_url:", v.Req_url, " des_url:", v.Des_url)
			}
		case "help":
			fmt.Println("Example: add|del tag rurl durl")
		case "quit":
			cancel()
			for k := range cycletaskmap {
				cycletaskmap[k].Wait()
			}
			SaveConfig()
			fmt.Println("Thanks!")
			wg.Done()
			return
		}
	}
}

type CycleTaskInfo struct {
	Tag          string `json:"tag"`
	Req_url      string `json:"rurl"`
	Des_url      string `json:"durl"`
	TimeInterval int    `json:"time"`
}

type CycleTask struct {
	sync.WaitGroup
	CycleTaskInfo
	Ticker      *time.Ticker
	Stop_signal chan bool
}

func NewCycleTask(r string, d string, tag string, t int) (c *CycleTask) {
	c = new(CycleTask)
	c.Stop_signal = make(chan bool)
	c.Req_url = r
	c.Des_url = d
	c.Tag = tag
	c.TimeInterval = t
	c.Ticker = time.NewTicker(time.Second * time.Duration(c.TimeInterval))
	return
}

func (c *CycleTask) StopCycle() {
	c.Stop_signal <- true
}

func (c *CycleTask) StartCycle(ctx context.Context) {
	c.Add(1)
	go func() {
		for {
			select {
			case <-c.Ticker.C:
				req_res, err := http.Get(c.Req_url)
				if err != nil {
					continue
				}
				// data, _ := ioutil.ReadAll(res.Body)
				// body := bytes.NewReader(data)
				_, err = http.Post(c.Des_url, "multipart/form-data", req_res.Body)
				defer req_res.Body.Close()
				if err != nil {
					continue
				}
			case <-c.Stop_signal:
				c.Ticker.Stop()
				c.Done()
				return
			case <-ctx.Done():
				c.Ticker.Stop()
				c.Done()
				return
			}
		}
	}()
}

func LoadConfig() {
	cycletaskmap = make(map[string]*CycleTask)

	filename := "config.json"
	var ct []CycleTaskInfo

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("没有发现config.json文件!")
		return
	}
	err = json.Unmarshal(data, &ct)
	if err != nil {
		fmt.Println("解析config.json文件失败!")
		return
	}
	for _, v := range ct {
		cycletask := NewCycleTask(v.Req_url, v.Des_url, v.Tag, v.TimeInterval)
		cycletaskmap[v.Tag] = cycletask
	}
}

func SaveConfig() {
	var infolist []CycleTaskInfo
	for _, v := range cycletaskmap {
		var info CycleTaskInfo
		info.Des_url = v.Des_url
		info.Req_url = v.Req_url
		info.Tag = v.Tag
		info.TimeInterval = v.TimeInterval
		infolist = append(infolist, info)
	}
	data, err := json.Marshal(infolist)
	if err != nil {
		fmt.Println("编码成json文件出错!")
	}
	err = ioutil.WriteFile("config.json", data, 0666)
	if err != nil {
		fmt.Println("保存config.json文件出错!")
	}
}

func main() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill)
	LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go HandleInput(ctx, cancel)

	wg.Wait()
}
