# Overview

The city information is from [addrparser](https://github.com/hsp8712/addrparser), it includes china's all province/city/region's code, center points, polyline points. It's based on GCJ-02. The Point-In-Polygon Algorithm is from [golang-geo](https://github.com/kellydunn/golang-geo)


## Install

```
$ go get -t github.com/vagrant-tiger/geo2city
```

## 下载数据文件
[china-region-20190902.zip](https://github.com/hsp8712/addrparser/releases/download/addrparser-1.0.1/china-region-20190902.zip) 解压后得到文件：china-region.json，文件中包括了所有的省市区行政区域信息，包括编码、名称、中心点、边界点集合。


## Usage Examples

```go
package main

import (
	"fmt"
	"github.com/vagrant-tiger/geo2city"
)


func main() {
	// 初始化数据并执行位置解析引擎，只有省的数据的话，大约300ms，如果全量数据在3s左右
	e, err := geo2city.LocationParseEngin("path/china-region.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 解析确定位置
	l := e.Parse(118.750934, 32.038634)
	
	// 获取省市区信息
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
```