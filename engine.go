package geo2city

import (
	"encoding/json"
	"errors"
	"github.com/kellydunn/golang-geo"
	"strconv"
	"strings"
)

const (
	COUNTRY  = iota
	PROVINCE = iota
	CITY     = iota
	DISTRICT = iota
	STREET   = iota
)

var locationEngin *LocationEngin

func init() {
	locationEngin = NewLocationEngin()
}

func LocationParseEngin() (*LocationParserEngine, error) {
	path := "resource/china-region.json"
	if le, ok := locationEngin.Load(path); ok {
		return le, nil
	}

	var err error
	rdByte, err := Asset(path)
	if err != nil {
		return nil, errors.New("can't find file: " + path + " " + err.Error())
	}

	rd := make([]resourceData, 0)
	err = json.Unmarshal(rdByte, &rd)
	if err != nil {
		return nil, errors.New("unmarshal resourceData failed:" + err.Error())
	}

	le := &LocationParserEngine{
		ProvinceMap: make(map[uint]*RegionInfo),
		CityMap:     make(map[uint]map[uint]*RegionInfo),
		RegionMap:   make(map[uint]map[uint]*RegionInfo),
	}
	for _, r := range rd {
		if r.Level == PROVINCE {
			l, err := r.convert()
			if err != nil {
				return nil, err
			}
			le.ProvinceMap[r.Code] = l
		} else if r.Level == CITY {
			ct, ok := le.CityMap[r.ParentCode]
			if !ok {
				ct = make(map[uint]*RegionInfo)
				le.CityMap[r.ParentCode] = ct
			}
			l, err := r.convert()
			if err != nil {
				return nil, err
			}
			ct[r.Code] = l
		} else if r.Level == DISTRICT || r.Level == STREET {
			rg, ok := le.RegionMap[r.ParentCode]
			if !ok {
				rg = make(map[uint]*RegionInfo)
				le.RegionMap[r.ParentCode] = rg
			}
			l, err := r.convert()
			if err != nil {
				return nil, err
			}
			rg[r.Code] = l
		}
	}

	locationEngin.Store(path, le)

	return le, nil
}

type resourceData struct {
	Name       string `json:"name"`
	Level      uint   `json:"level"`
	Code       uint   `json:"code"`
	ParentCode uint   `json:"parentCode"`
	Center     string `json:"center"`
	Polyline   string `json:"polyline"`
}

func (rd resourceData) convert() (l *RegionInfo, err error) {
	centerPoint, err := toPoint(rd.Center)
	if err != nil {
		return l, err
	}
	s := strings.Split(rd.Polyline, "|")
	polyline := make([]Polygon, 0)
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

	l = &RegionInfo{
		Name:       rd.Name,
		Level:      rd.Level,
		Code:       rd.Code,
		ParentCode: rd.ParentCode,
		Center:     centerPoint,
		Polyline:   polyline,
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
	ProvinceMap map[uint]*RegionInfo
	CityMap     map[uint]map[uint]*RegionInfo
	RegionMap   map[uint]map[uint]*RegionInfo
}

func (lpe *LocationParserEngine) Parse(lat float64, lng float64) (location Location) {
	location = Location{}
	searchPoint := geo.NewPoint(lat, lng)

	// search in province
	for _, pr := range lpe.ProvinceMap {
		if pr.contain(searchPoint) {
			location.prov = pr
			break
		}
	}

	if location.prov == nil || len(lpe.CityMap) == 0 {
		return location
	}

	// search in city
	cities := lpe.CityMap[location.prov.Code]
	for _, cr := range cities {
		if cr.contain(searchPoint) {
			location.city = cr
			break
		}
	}

	if location.city == nil || len(lpe.RegionMap) == 0 {
		return location
	}

	// search in district
	districts := lpe.RegionMap[location.city.Code]
	for _, dr := range districts {
		if dr.contain(searchPoint) {
			location.district = dr
			break
		}
	}

	return location
}

func (l Location) GetProv() (*RegionInfo, error) {
	if l.prov == nil {
		return l.prov, errors.New("province not found")
	}
	return l.prov, nil
}

func (l Location) GetCity() (*RegionInfo, error) {
	if l.city == nil {
		return l.city, errors.New("city not found")
	}
	return l.city, nil
}

func (l Location) GetDistrict() (*RegionInfo, error) {
	if l.district == nil {
		return l.district, errors.New("district not found")
	}
	return l.district, nil
}

type Location struct {
	prov     *RegionInfo
	city     *RegionInfo
	district *RegionInfo
}

type RegionInfo struct {
	Name       string     `json:"name"`
	Level      uint       `json:"level"`
	Code       uint       `json:"code"`
	ParentCode uint       `json:"parentCode"`
	Center     *geo.Point `json:"center"`
	Polyline   []Polygon  `json:"polyline"`
}

func (ri *RegionInfo) contain(point *geo.Point) bool {
	for _, pl := range ri.Polyline {
		if pl.contain(point) {
			return true
		}
	}
	return false
}

func (ri *RegionInfo) getName() string {
	return ri.Name
}

type Polygon []*geo.Point

func (p Polygon) contain(point *geo.Point) bool {
	polygonGeo := geo.NewPolygon(p)
	return polygonGeo.Contains(point)
}
