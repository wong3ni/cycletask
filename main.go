package main

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
)

type Code int

type program struct {
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() error {
	cfg = new(Config)
	cfg.Load()
	if cfg.ListenAddr == "" {
		cfg.Save()
		fmt.Println("已生成空配置文件，请设置监听地址!")
		return fmt.Errorf("init failed!")
	}
	logger = new(Logger)
	logger.Load()
	logger.Start()
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
		Description: "cycletask demo",
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
		}
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
