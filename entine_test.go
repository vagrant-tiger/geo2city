package geo2city

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLocationEngin(t *testing.T) {
	startTime := time.Now()
	e, err := LocationParseEngin()
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
		fmt.Println(prov.getName())
	} else {
		fmt.Println(err.Error())
	}
	city, err := l.GetCity()
	if err == nil {
		fmt.Println(city.getName())
	} else {
		fmt.Println(err.Error())
	}
	district, err := l.GetDistrict()
	if err == nil {
		fmt.Println(district.getName())
	} else {
		fmt.Println(err.Error())
	}
}
