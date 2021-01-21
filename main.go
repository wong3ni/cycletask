package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

type Code int

var ResPath string

type program struct {
}

func (p *program) Start(s service.Service) error {
	NewYolo()
	SetWeights("/Users/wzh/model/yolov4_final.weights")
	SetCfg("/Users/wzh/Python3/image_analysis/cfg/yolov4.cfg")
	SetCategory("/Users/wzh/Python3/image_analysis/data/coco.names")
	if service.Interactive() == false {
		path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		ResPath = path
	} else {
		SetPrintLayers(1)
		SetPrintDetectTime(1)
	}
	InitYoloNet()
	go p.run()
	return nil
}

func (p *program) run() error {
	cfg = new(Config)
	cfg.Load()
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":22222"
	}
	if cfg.Path == "" {
		cfg.Path = "logs"
		cfg.MaxLines = 30000
	}
	logger = new(Logger)
	logger.Load()
	if service.Interactive() == false {
		logger.Println = logger.ToFile
		logger.Start()
	} else {
		logger.Println = logger.ToShow
	}
	CyT = new(CycleTask)
	CyT.Load()
	htsv = new(HTTPServer)
	htsv.StartHTTPServer()
	return nil
}

func (p *program) Stop(s service.Service) error {
	CyT.Stop()
	CyT.Save()
	htsv.srv.Close()
	// htsv.Wait()
	logger.Close()
	cfg.Save()
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "cycletask",
		DisplayName: "cycletask",
		Description: "cycletask daemon",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = s.Install()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if os.Args[1] == "uninstall" {
			err = s.Uninstall()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if os.Args[1] == "stop" {
			err = s.Stop()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if os.Args[1] == "start" {
			err = s.Start()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
