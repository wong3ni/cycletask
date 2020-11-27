package main

import (
	"context"
	"fmt"
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
			cycletask := NewCycleTask(req_url, des_url)
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
			for k := range cycletaskmap {
				fmt.Println("tag:", k, " req_url:", cycletaskmap[k].Req_url)
			}
		case "help":
			fmt.Println("Example: add|del tag url")
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

func NewCycleTask(r string, d string) (c *CycleTask) {
	c = new(CycleTask)
	c.Stop_signal = make(chan bool)
	c.Req_url = r
	c.Des_url = d
	c.TimeInterval = 2
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

}

func SaveConfig() {
	for tag := range cycletaskmap {
		fmt.Println("tag:", tag, " req_url:", cycletaskmap[tag].Req_url)
	}
}

func main() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill)
	cycletaskmap = make(map[string]*CycleTask)

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go HandleInput(ctx, cancel)

	wg.Wait()
}
