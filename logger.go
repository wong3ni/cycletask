package main

import (
	"fmt"
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
	fd *os.File
}

func (l *Logger) GetFileName() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d %d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) + ".txt"
}

func (l *Logger) Start() {
	// abspath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.MkdirAll(filepath.Join(ResPath, l.Path), os.ModePerm)
	l.fd, _ = os.Create(filepath.Join(ResPath, l.Path, l.GetFileName()))
	log.SetOutput(l.fd)
}

func (l *Logger) Close() {
	if l.fd != nil {
		l.fd.Close()
	}
}

func (l *Logger) Load() {
	l.LoggerInfo = cfg.LoggerInfo
}

var logger *Logger
