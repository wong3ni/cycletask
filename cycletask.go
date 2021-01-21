package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CycleTaskUnitInfo struct {
	Tag          string `json:"tag"`
	Req_url      string `json:"rurl"`
	Des_url      string `json:"durl"`
	TimeInterval int    `json:"time"`
	State        bool   `json:"state"`
	Direction    string `json:"direction"`
	Name         string `json:"name"`
	Id           string `json:"id"`
	NVRID        string `json:"nvrid"`
	Channel      int    `json:"channel"`
}

type CycleTask struct {
	Cycletaskmap sync.Map
}

var CyT *CycleTask

func (c *CycleTask) Load() {
	for _, v := range cfg.CycleTaskUnitInfos {
		cyctu := NewCycleTaskUnit(v)
		c.Cycletaskmap.Store(v.Tag, cyctu)
	}
}

func (c *CycleTask) Save() {
	cfg.CycleTaskUnitInfos = cfg.CycleTaskUnitInfos[0:0]
	c.Cycletaskmap.Range(func(k, v interface{}) bool {
		cyctu := v.(*CycleTaskUnit)
		cfg.CycleTaskUnitInfos = append(cfg.CycleTaskUnitInfos, cyctu.CycleTaskUnitInfo)
		return true
	})
}

func (c *CycleTask) Stop() {
	c.Cycletaskmap.Range(func(k, v interface{}) bool {
		cyctu := v.(*CycleTaskUnit)
		if cyctu.State {
			cyctu.StopCycle()
			cyctu.Wait()
		}
		return true
	})
}

func (c *CycleTask) GetTaskUnit(tag string) (*CycleTaskUnit, bool) {
	cyctu, ok := c.Cycletaskmap.Load(tag)
	if cyctu != nil {
		return cyctu.(*CycleTaskUnit), ok
	}
	return nil, false
}

func (c *CycleTask) AddTaskUnit(cyctuinfo CycleTaskUnitInfo) Code {
	cyctu := NewCycleTaskUnit(cyctuinfo)
	_, ok := c.Cycletaskmap.LoadOrStore(cyctuinfo.Tag, cyctu)
	if ok {
		return 1
	}
	return 0
}

func (c *CycleTask) StartTaskUnit(tag string) Code {
	cyctu, ok := c.GetTaskUnit(tag)
	if ok {
		go cyctu.StartCycle()
		return 0
	}
	return 1
}

func (c *CycleTask) StopTaskUnit(tag string) Code {
	cyctu, ok := c.GetTaskUnit(tag)
	if ok {
		cyctu.StopCycle()
		return 0
	}
	return 1
}

func (c *CycleTask) DelTaskUnit(tag string) Code {
	cyctu, ok := c.GetTaskUnit(tag)
	if ok {
		if cyctu.State {
			cyctu.StopCycle()
			cyctu.Wait()
		}
		c.Cycletaskmap.Delete(tag)
		return 0
	}
	return 1
}

type CycleTaskUnit struct {
	sync.WaitGroup
	CycleTaskUnitInfo
	Ticker      *time.Ticker
	stop_signal chan bool
}

func NewCycleTaskUnit(cyctuinfo CycleTaskUnitInfo) (c *CycleTaskUnit) {
	c = new(CycleTaskUnit)
	c.stop_signal = make(chan bool)
	c.CycleTaskUnitInfo = cyctuinfo
	return
}

func (c *CycleTaskUnit) StopCycle() {
	c.stop_signal <- true
}

func (c *CycleTaskUnit) StartCycle() {
	c.Ticker = time.NewTicker(time.Second * time.Duration(c.TimeInterval))
	c.Add(1)
	c.State = true
	for {
		select {
		case <-c.Ticker.C:
			go c.Deal()
		case <-c.stop_signal:
			c.Ticker.Stop()
			c.Ticker = nil
			c.State = false
			c.Done()
			return
		}
	}
}

func (c *CycleTaskUnit) Deal() {
	req, err := http.NewRequest("GET", c.Req_url, nil)
	req.Close = true
	myclient := http.Client{Timeout: 2 * time.Second}
	req_res, err := myclient.Do(req)
	if err != nil {
		logger.Println("Tag", c.Tag, err)
		return
	}
	defer req_res.Body.Close()

	data, _ := ioutil.ReadAll(req_res.Body)
	if len(data) < 15 {
		// logger.Println("Tag:", c.Tag, " <= ", "no such stream")
		if c.NVRID != "" {
			r_s := "http://" + strings.Split(c.Req_url, "/")[2] + "/gb28181/invite?id=" + c.NVRID + "&channel=" + fmt.Sprintf("%d", c.Channel)
			req, err = http.NewRequest("GET", r_s, nil)
			req.Close = true
			r_r, err := myclient.Do(req)
			if err != nil {
				logger.Println("Tag", c.Tag, err)
				return
			}
			if r_r != nil {
				logger.Println(r_r.StatusCode, r_s)
			}
		}
		return
	}
	r := DealImage(data, c.Tag)
	fmt.Println(r)
}
