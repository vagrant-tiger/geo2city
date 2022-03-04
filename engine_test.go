package geo2city

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLocationEngin(t *testing.T) {
	startTime := time.Now()
	e, err := LocationParseEngin("D:/git/geo2city/china-region.json")
	endTime := time.Now()
	useTime := endTime.Sub(startTime).Milliseconds()
	fmt.Printf("start use %d ms\n", useTime)
	if err != nil {
		t.Errorf("failed " + err.Error())
	}

	startTime = time.Now()
	l := e.Parse(118.750934, 32.038634)
	endTime = time.Now()
	useTime = endTime.Sub(startTime).Milliseconds()
	fmt.Printf("parse use %d ms\n", useTime)
	prov, err := l.GetProv()
	if err == nil {
		fmt.Println(prov.GetName())
	} else {
		fmt.Println(err.Error())
	}
	city, err := l.GetCity()
	if err == nil {
		fmt.Println(city.GetName())
	} else {
		fmt.Println(err.Error())
	}
	district, err := l.GetDistrict()
	if err == nil {
		fmt.Println(district.GetName())
	} else {
		fmt.Println(err.Error())
	}
}

func TestFlush2DB(t *testing.T) {
	err := Flush2DB("D:/git/geo2city/china-region.json", "127.0.0.1", "3306", "root", "123456", "china_geo")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestLocationDbEngin(t *testing.T) {
	startTime := time.Now()
	e, err := LocationDbEngin("127.0.0.1", "3306", "root", "123456", "china_geo")
	endTime := time.Now()
	useTime := endTime.Sub(startTime).Milliseconds()
	fmt.Printf("start use %d ms\n", useTime)
	if err != nil {
		t.Errorf("failed " + err.Error())
	}

	startTime = time.Now()
	l := e.Parse(118.750934, 32.038634)
	endTime = time.Now()
	useTime = endTime.Sub(startTime).Milliseconds()
	fmt.Printf("parse use %d ms\n", useTime)
	prov, err := l.GetProv()
	if err == nil {
		fmt.Println(prov.GetName())
	} else {
		fmt.Println(err.Error())
	}
	city, err := l.GetCity()
	if err == nil {
		fmt.Println(city.GetName())
	} else {
		fmt.Println(err.Error())
	}
	district, err := l.GetDistrict()
	if err == nil {
		fmt.Println(district.GetName())
	} else {
		fmt.Println(err.Error())
	}
}
