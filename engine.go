package geo2city

import (
	"encoding/json"
	"errors"
	"github.com/vagrant-tiger/geo2city/geo"
	"strconv"
	"strings"
	"sync"
)

const (
	COUNTRY  = iota
	PROVINCE = iota
	CITY     = iota
	DISTRICT = iota
	STREET   = iota
)

var (
	wg sync.WaitGroup
	le *locationEngin
)

func init() {
	le = newLocationEngin()
}

func LocationParseEngin() (*LocationParserEngine, error) {
	path := "resource/china-region.json"
	if lpe, ok := le.Load(path); ok {
		return lpe, nil
	}

	var err error
	rdByte, err := Asset(path)
	if err != nil {
		return nil, errors.New("can't find file: " + path + " " + err.Error())
	}

	rd := make([]*resourceData, 0)
	err = json.Unmarshal(rdByte, &rd)
	if err != nil {
		return nil, errors.New("unmarshal resourceData failed:" + err.Error())
	}

	lpe := &LocationParserEngine{
		provinceMap: make(map[uint]*regionInfo),
		cityMap:     make(map[uint]map[uint]*regionInfo),
		regionMap:   make(map[uint]map[uint]*regionInfo),
	}

	regionChan := make(chan *regionInfo, 100)
	go storeRegion(regionChan, lpe)

	for _, r := range rd {
		wg.Add(1)
		go sendRegion(regionChan, r)
	}
	wg.Wait()

	le.Store(path, lpe)

	return lpe, nil
}

func sendRegion(c chan *regionInfo, r *resourceData) {
	l, _ := r.convert()
	c <- l
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

type resourceData struct {
	Name       string `json:"name"`
	Level      uint   `json:"level"`
	Code       uint   `json:"code"`
	ParentCode uint   `json:"parentCode"`
	Center     string `json:"center"`
	Polyline   string `json:"polyline"`
}

func (rd resourceData) convert() (l *regionInfo, err error) {
	centerPoint, err := toPoint(rd.Center)
	if err != nil {
		return l, err
	}
	s := strings.Split(rd.Polyline, "|")
	polyline := make([]polygon, 0)
	for _, ps := range s {
		pg := make([]*geo.Point, 0)
		pgs := strings.Split(ps, ";")
		for _, pgsp := range pgs {
			p, err := toPoint(pgsp)
			if err != nil {
				return l, err
			}
			pg = append(pg, p)
		}
		polyline = append(polyline, pg)
	}

	l = &regionInfo{
		name:       rd.Name,
		level:      rd.Level,
		code:       rd.Code,
		parentCode: rd.ParentCode,
		center:     centerPoint,
		polyline:   polyline,
	}
	return l, nil
}

func toPoint(data string) (p *geo.Point, err error) {
	l := strings.Split(data, ",")
	lat, err := strconv.ParseFloat(l[0], 64)
	if err != nil {
		return p, errors.New("parse center data failed:" + data)
	}
	lng, err := strconv.ParseFloat(l[1], 64)
	if err != nil {
		return p, errors.New("parse center data failed:" + data)
	}

	p = geo.NewPoint(lat, lng)
	return p, nil
}

type LocationParserEngine struct {
	provinceMap map[uint]*regionInfo
	cityMap     map[uint]map[uint]*regionInfo
	regionMap   map[uint]map[uint]*regionInfo
}

func (lpe *LocationParserEngine) Parse(lat float64, lng float64) (lc location) {
	lc = location{}
	searchPoint := geo.NewPoint(lat, lng)

	// search in province
	for _, pr := range lpe.provinceMap {
		if pr.contain(searchPoint) {
			lc.prov = pr
			break
		}
	}

	if lc.prov == nil || len(lpe.cityMap) == 0 {
		return lc
	}

	// search in city
	cities := lpe.cityMap[lc.prov.code]
	for _, cr := range cities {
		if cr.contain(searchPoint) {
			lc.city = cr
			break
		}
	}

	if lc.city == nil || len(lpe.regionMap) == 0 {
		return lc
	}

	// search in district
	districts := lpe.regionMap[lc.city.code]
	for _, dr := range districts {
		if dr.contain(searchPoint) {
			lc.district = dr
			break
		}
	}

	return lc
}

func (l location) GetProv() (*regionInfo, error) {
	if l.prov == nil {
		return l.prov, errors.New("province not found")
	}
	return l.prov, nil
}

func (l location) GetCity() (*regionInfo, error) {
	if l.city == nil {
		return l.city, errors.New("city not found")
	}
	return l.city, nil
}

func (l location) GetDistrict() (*regionInfo, error) {
	if l.district == nil {
		return l.district, errors.New("district not found")
	}
	return l.district, nil
}

type location struct {
	prov     *regionInfo
	city     *regionInfo
	district *regionInfo
}

type regionInfo struct {
	name       string     `json:"name"`
	level      uint       `json:"level"`
	code       uint       `json:"code"`
	parentCode uint       `json:"parentCode"`
	center     *geo.Point `json:"center"`
	polyline   []polygon  `json:"polyline"`
}

func (ri *regionInfo) contain(point *geo.Point) bool {
	for _, pl := range ri.polyline {
		if pl.contain(point) {
			return true
		}
	}
	return false
}

func (ri *regionInfo) GetName() string {
	return ri.name
}

type polygon []*geo.Point

func (p polygon) contain(point *geo.Point) bool {
	polygonGeo := geo.NewPolygon(p)
	return polygonGeo.Contains(point)
}
