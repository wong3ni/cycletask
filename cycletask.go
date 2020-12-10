package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
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
}

type CycleTaskUnit struct {
	sync.WaitGroup
	CycleTaskUnitInfo
	Ticker      *time.Ticker
	stop_signal chan bool
}

type CycleTask struct {
	Cycletaskmap sync.Map
}

var CyT *CycleTask

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
	go func() {
		for {
			select {
			case <-c.Ticker.C:
				req_res, err := http.Get(c.Req_url)
				if req_res != nil {
					log.Println("Tag:", c.Tag, " <= ", req_res.Status, c.Req_url)
				} else {
					log.Println("Tag:", c.Tag, " <= ", "unreachable", c.Req_url)
				}
				if err != nil {
					continue
				}
				data, _ := ioutil.ReadAll(req_res.Body)
				if len(data) < 15 {
					log.Println("Tag:", c.Tag, " <= ", "no such stream")
					// 	r_s := "http://" + strings.Split(c.Req_url, "/")[2] + "/gb28181/invite?id=" + c.NVRID + "&channel=0"
					// 	http.Get(r_s)
					// 	log.Println(r_s)
					continue
				}
				body := bytes.NewReader(data)
				req, err := http.NewRequest("POST", c.Des_url, body)
				req.Header.Set("contentType", "multipart/form-data")
				req.Header.Set("direction", c.Direction)
				req.Header.Set("name", c.Name)
				req.Header.Set("id", c.Id)
				req.Header.Set("tag", c.Tag)
				res_res, err := http.DefaultClient.Do(req)

				// res_res, err := http.Post(c.Des_url, "multipart/form-data", req_res.Body)
				if res_res != nil {
					log.Println("Tag:", c.Tag, " => ", res_res.Status, c.Des_url)
				} else {
					log.Println("Tag:", c.Tag, " => ", "unreachable", c.Des_url)
				}
				defer req_res.Body.Close()
				if err != nil {
					continue
				}
			case <-c.stop_signal:
				c.Ticker.Stop()
				c.Ticker = nil
				c.Done()
				return
			}
		}
	}()
}

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
			cyctu.State = false
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
		cyctu.StartCycle()
		cyctu.State = true
		return 0
	}
	return 1
}

func (c *CycleTask) StopTaskUnit(tag string) Code {
	cyctu, ok := c.GetTaskUnit(tag)
	if ok {
		cyctu.StopCycle()
		cyctu.State = false
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

// func SendMonibucaInvite(id, channel string) bool {
// 	u := fmt.Sprintf("http://219.138.126.226:18298/gb28181/invite?id=%s&channel=%s", id, channel)
// 	res, _ := http.Get(u)
// 	if res != nil {
// 		if res.StatusCode == 200 {
// 			return true
// 		}
// 	}
// 	return false
// }
