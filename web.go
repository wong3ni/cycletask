package main

import (
	"encoding/json"
	"html/template"
	"log"
	"sort"
	"strconv"

	"net/http"
	"time"

	"sync"
)

type WebConfig struct {
	ListenAddr string `json:"ListenAddr"`
}

type HTTPServer struct {
	sync.WaitGroup
	srv *http.Server
}

type Res struct {
	Cod Code   `json:"Code"`
	Msg string `json:"Msg"`
}

func NewRes() *Res {
	res := new(Res)
	res.Cod = -1
	res.Msg = "unknown"
	return res
}

func MiddleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (h *HTTPServer) StartHTTPServer() {
	h.srv = &http.Server{Addr: cfg.ListenAddr}

	http.Handle("/static/", MiddleLog(http.StripPrefix("/static/", http.FileServer(http.Dir("dist/static")))))

	http.Handle("/", MiddleLog(http.HandlerFunc(index)))
	http.Handle("/api/add", MiddleLog(http.HandlerFunc(ApiCycleTaskAdd)))
	http.Handle("/api/start", MiddleLog(http.HandlerFunc(ApiCycleTaskStart)))
	http.Handle("/api/stop", MiddleLog(http.HandlerFunc(ApiCycleTaskStop)))
	http.Handle("/api/del", MiddleLog(http.HandlerFunc(ApiCycleTaskDel)))
	http.Handle("/cycletask/list", http.HandlerFunc(ApiCycleTaskList))
	http.Handle("/api/test", MiddleLog(http.HandlerFunc(ApiTest)))

	h.Add(1)
	log.Printf("HTTP Server Start: %s \n", cfg.ListenAddr)
	go func() {
		defer h.Done()
		if err := h.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("检测监听地址是否有误!")
			log.Fatalln(err)
		}
	}()
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("dist/index.html")
	t.Execute(w, nil)
}

func ApiCycleTaskAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := NewRes()
	if r.Method == "GET" {
		r.ParseForm()
		var cyctuinfo CycleTaskUnitInfo
		cyctuinfo.Req_url = r.Form.Get("rurl")
		cyctuinfo.Des_url = r.Form.Get("durl")
		cyctuinfo.Tag = r.Form.Get("tag")
		cyctuinfo.Direction = r.Form.Get("direction")
		cyctuinfo.Name = r.Form.Get("name")
		cyctuinfo.Id = r.Form.Get("id")
		cyctuinfo.NVRID = r.Form.Get("nvrid")
		t := r.Form.Get("time")
		ti, err := strconv.Atoi(t)
		if err != nil || ti <= 0 {
			ti = 2
		}
		cyctuinfo.TimeInterval = ti
		res.Cod = CyT.AddTaskUnit(cyctuinfo)
		if res.Cod == 0 {
			res.Msg = "success"
		} else {
			res.Msg = "error"
		}
		resjson, _ := json.Marshal(res)
		w.Write(resjson)
	}
}

func ApiCycleTaskStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := NewRes()
	if r.Method == "GET" {
		r.ParseForm()
		tag := r.Form["tag"][0]
		res.Cod = CyT.StartTaskUnit(tag)
		res.Msg = "success"
		resjson, _ := json.Marshal(res)
		w.Write(resjson)
	}
}

func ApiCycleTaskStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := NewRes()
	if r.Method == "GET" {
		r.ParseForm()
		tag := r.Form["tag"][0]
		res.Cod = CyT.StopTaskUnit(tag)
		res.Msg = "success"
		resjson, _ := json.Marshal(res)
		w.Write(resjson)
	}
}

func ApiCycleTaskDel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := NewRes()
	if r.Method == "GET" {
		r.ParseForm()
		tag := r.Form["tag"][0]
		res.Cod = CyT.DelTaskUnit(tag)
		res.Msg = "success"
		resjson, _ := json.Marshal(res)
		w.Write(resjson)
	}
}

func ApiCycleTaskList(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "Start EventSource", r.RequestURI)
	sse := NewSSE(w, r.Context())
	for {
		var list []CycleTaskUnitInfo
		CyT.Cycletaskmap.Range(func(k, v interface{}) bool {
			cyctu := v.(*CycleTaskUnit)
			list = append(list, cyctu.CycleTaskUnitInfo)
			return true
		})
		sort.Slice(list, func(i, j int) bool {
			return list[i].Tag < list[j].Tag
		})

		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		sse.WriteJSON(list)
		select {
		case <-ticker.C:
			if sse.WriteJSON(list) != nil {
				return
			}
		case <-r.Context().Done():
			log.Println(r.RemoteAddr, "Close EventSource", r.RequestURI)
			return
		}
	}
}

func ApiTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := NewRes()
	if r.Method == "GET" {
		r.ParseForm()
		resjson, _ := json.Marshal(res)
		w.Write(resjson)
	}
}
