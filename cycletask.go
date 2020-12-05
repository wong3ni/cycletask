package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CycleTaskUnitInfo struct {
	Tag          string `json:"tag"`
	Req_url      string `json:"rurl"`
	Des_url      string `json:"durl"`
	TimeInterval int    `json:"time"`
	State        bool   `json:"state"`
	Direction    int    `json:"direction"`
	Name		 string  `json:"name"`
	Id   	     string  `json:id`

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

func NewCycleTaskUnit(r string, d string, tag string, t int, dire int, name string, id string) (c *CycleTaskUnit) {
	c = new(CycleTaskUnit)
	c.stop_signal = make(chan bool)
	c.Req_url = r
	c.Des_url = d
	c.Tag = tag
	c.TimeInterval = t
	c.State = false
	c.Direction = dire
	c.Name = name
	c.Id = id
	// c.Ticker = time.NewTicker(time.Second * time.Duration(c.TimeInterval))
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
				// data, _ := ioutil.ReadAll(res.Body)
				// body := bytes.NewReader(data)
				req, err := http.NewRequest("POST", c.Des_url, req_res.Body)
				req.Header.Set("contentType", "multipart/form-data")
				req.Header.Set("direction", strconv.Itoa(c.Direction))
				req.Header.Set("name", c.Name)
				req.Header.Set("id", c.Id)
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
		cyctu := NewCycleTaskUnit(v.Req_url, v.Des_url, v.Tag, v.TimeInterval, v.Direction, v.Name, v.Id)
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

func (c *CycleTask) AddTaskUnit(req_url, des_url, tag string, timeinterval , direction int, name string, id string) Code {
	cyctu := NewCycleTaskUnit(req_url, des_url, tag, timeinterval, direction, name, id)
	_, ok := c.Cycletaskmap.LoadOrStore(tag, cyctu)
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
