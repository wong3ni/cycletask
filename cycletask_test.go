package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpPost(t *testing.T) {
	des_url := "http://192.168.168.100:8080/boat"
	byte, err := ioutil.ReadFile("/Users/wzh/GoProject/vidprocessing/core/frame.data")
	if err != nil {
		fmt.Println("读取文件失败", err)
		return
	}
	res, err := http.Post(des_url, "multipart/form-data", bytes.NewReader(byte))
	if err != nil {
		fmt.Println("请求失败", err)
		return
	}
	fmt.Println(res)
}

type Website struct {
	Name   string `xml:"name,attr"`
	Url    string
	Course []string
}

func TestRWJsonFile(t *testing.T) {

}
