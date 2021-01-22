package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -L ./ -lshipcount
#cgo CXXFLAGS: -I ./

void newYolo();
void deleteYolo();
void initYoloNet();

void setWeights(const char *);
void setCfg(const char *);
void setCategory(const char *);

void setPrintLayers(int);
void setCuda(int);
void setPrintDetectTime(int);

unsigned char deal(unsigned char *, const char *);
*/
import "C"

type YoloConfig struct {
	Weights  string `json:"weights"`
	Cfg      string `json:"cfg"`
	Category string `json:"category"`
}

func NewYolo() {
	C.newYolo()
}

func DeleteYolo() {
	C.deleteYolo()
}

func InitYoloNet() {
	C.initYoloNet()
}

func SetWeights(s string) {
	C.setWeights(C.CString(s))
}

func SetCfg(s string) {
	C.setCfg(C.CString(s))
}

func SetCategory(s string) {
	C.setCategory(C.CString(s))
}

func SetPrintLayers(i int) {
	C.setPrintLayers(C.int(i))
}

func SetCuda(i int) {
	C.setCuda(C.int(i))
}

func SetPrintDetectTime(i int) {
	C.setPrintDetectTime(C.int(i))
}

func DealImage(data []byte, s string) uint8 {
	return byte(C.deal((*C.uchar)(C.CBytes(data)), C.CString(s)))
}
