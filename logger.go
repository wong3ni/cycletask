package main

import (
	"log"
	"os"
	"path/filepath"

	"time"
)

type LoggerInfo struct {
	Path string `json:"path"`
}

type Logger struct {
	LoggerInfo
	fd     *os.File
	buffer chan string
	stop   chan bool
}

func (l *Logger) GetFileName() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05") + ".txt"
}

func (l *Logger) Start() {
	l.buffer = make(chan string, 10)
	l.stop = make(chan bool)
	os.MkdirAll(l.Path, os.ModePerm)
	l.fd, _ = os.Create(filepath.Join(l.Path, l.GetFileName()))
	log.SetOutput(l.fd)
	go l.loop()
}

func (l *Logger) loop() {
	for {
		select {
		case s := <-l.buffer:
			log.Println(s)
		case <-l.stop:
			return
		}
	}
}

func (l *Logger) Close() {
	l.fd.Close()
	l.stop <- true
}

func (l *Logger) Write(s string) {
	l.buffer <- s
}

func (l *Logger) Load() {
	l.LoggerInfo = cfg.LoggerInfo
}

var logger *Logger
