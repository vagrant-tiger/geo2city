# Overview

The city information is from [addrparser](https://github.com/hsp8712/addrparser), it includes china's all province/city/region's code, center points, polyline points. It's based on GCJ-02. The Point-In-Polygon Algorithm is from [golang-geo](https://github.com/kellydunn/golang-geo)


## Install

```
$ go get -t github.com/vagrant-tiger/geo2city
```

## download geo data
[china-region-20190902.zip](https://github.com/hsp8712/addrparser/releases/download/addrparser-1.0.1/china-region-20190902.zip)


## Usage Examples

### 1. use json file

```go
package main

import (
	"fmt"
	"github.com/vagrant-tiger/geo2city"
)


func main() {
	// 初始化数据并执行位置解析引擎
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

### 2. use mysql

#### fist: flush data to mysql

```go
package main

import (
	"fmt"
	"github.com/vagrant-tiger/geo2city"
)


func main() {
	// 数据导入mysql
	err := geo2city.Flush2DB("path/china-region.json", "host", "port", "user", "password", "database")
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

#### second: use mysql

```go
package main

import (
	"fmt"
	"github.com/vagrant-tiger/geo2city"
)


func main() {
	// 初始化数据并执行位置解析引擎
	e, err := geo2city.LocationDbEngin("host", "port", "user", "password", "database")
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