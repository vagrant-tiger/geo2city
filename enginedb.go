package geo2city

import (
	"encoding/json"
	"fmt"
	geo "github.com/kellydunn/golang-geo"
	"github.com/vagrant-tiger/geo2city/tool"
	"sync"
)

var wg sync.WaitGroup

func Flush2DB(path string, host string, port string, userName string, pass string, dbName string) error {
	// load data from file
	lpe, err := LocationParseEngin(path)
	if err != nil {
		return err
	}

	// create table
	fmt.Println("start create table:")
	sc := tool.NewMysqlConf(host, port, userName, pass, dbName)
	err = sc.CreateTable()
	if err != nil {
		return err
	}

	fmt.Println("start insert province data:")
	for _, pr := range lpe.provinceMap {
		err = pr.insert2db()
		if err != nil {
			return err
		}
	}
	fmt.Println("insert province finished")

	fmt.Println("start insert city data:")
	for _, v := range lpe.cityMap {
		for _, cr := range v {
			err = cr.insert2db()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("insert city finished")

	fmt.Println("start insert region data:")
	for _, v := range lpe.regionMap {
		for _, r := range v {
			err = r.insert2db()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ri *regionInfo) insert2db() error {
	centerPoint := point{
		Lat: ri.center.Lat(),
		Lng: ri.center.Lng(),
	}
	cp, _ := json.Marshal(centerPoint)
	pl, _ := json.Marshal(ri.polyline)
	err := tool.InsertChinaGeo(ri.code, ri.parentCode, ri.name, ri.level, string(cp), string(pl))
	if err != nil {
		return err
	}

	return nil
}

type point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type chinaGeo struct {
	id         uint   `json:"id"`
	name       string `json:"name"`
	level      uint   `json:"level"`
	code       uint   `json:"code"`
	parentCode uint   `json:"parentCode"`
	center     string `json:"center"`
	polyline   string `json:"polyline"`
}

func LocationDbEngin(host string, port string, userName string, pass string, dbName string) (*LocationParserEngine, error) {
	lpe := &LocationParserEngine{
		provinceMap: make(map[uint]*regionInfo),
		cityMap:     make(map[uint]map[uint]*regionInfo),
		regionMap:   make(map[uint]map[uint]*regionInfo),
	}

	sc := tool.NewMysqlConf(host, port, userName, pass, dbName)

	rows, err := sc.QueryData()
	if err != nil {
		return nil, err
	}

	regionChan := make(chan *regionInfo, 100)
	go storeRegion(regionChan, lpe)

	for rows.Next() {
		wg.Add(1)
		cg := &chinaGeo{}
		rows.Scan(&cg.id, &cg.code, &cg.parentCode, &cg.name, &cg.level, &cg.center, &cg.polyline)
		go sendRegion(regionChan, cg)
	}
	wg.Wait()

	return lpe, nil
}

func sendRegion(c chan *regionInfo, cg *chinaGeo) {
	var cp point
	json.Unmarshal([]byte(cg.center), &cp)

	centerGeo := geo.NewPoint(cp.Lat, cp.Lng)

	var pl []polygon

	json.Unmarshal([]byte(cg.polyline), &pl)
	ri := &regionInfo{
		name:       cg.name,
		level:      cg.level,
		code:       cg.code,
		parentCode: cg.parentCode,
		center:     centerGeo,
		polyline:   pl,
	}
	c <- ri
}

func storeRegion(c chan *regionInfo, lpe *LocationParserEngine) {
	for ri := range c {
		if ri.level == PROVINCE {
			lpe.provinceMap[ri.code] = ri
		} else if ri.level == CITY {
			ct, ok := lpe.cityMap[ri.parentCode]
			if !ok {
				ct = make(map[uint]*regionInfo)
				lpe.cityMap[ri.parentCode] = ct
			}
			ct[ri.code] = ri
		} else if ri.level == DISTRICT || ri.level == STREET {
			rg, ok := lpe.regionMap[ri.parentCode]
			if !ok {
				rg = make(map[uint]*regionInfo)
				lpe.regionMap[ri.parentCode] = rg
			}
			rg[ri.code] = ri
		}
		wg.Done()
	}
}
