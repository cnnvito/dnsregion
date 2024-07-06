package dnsregion

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

//go:embed data/ip2region.xdb
var database []byte

var DefaultIPDatabase = NewIPDatabase(database)

type IPParserRegion interface {
	Parser(ip string) IPResult
}
type IPResult struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
}

func (i IPResult) String() string {
	r, _ := json.Marshal(&i)
	return string(r)
}

func newEmptyIPResult(ip string) IPResult {
	return IPResult{
		IP:      ip,
		Country: "unknown",
		Region:  "unknown",
		City:    "unknown",
		ISP:     "unknown",
	}
}

func NewIPDatabase(db []byte) *IPDatabase {
	s, err := xdb.NewWithBuffer(db)
	if err != nil {
		panic(err)
	}
	return &IPDatabase{s}
}

type IPDatabase struct {
	*xdb.Searcher
}

func (s *IPDatabase) Parser(ip string) IPResult {
	result := newEmptyIPResult(ip)
	r, err := s.SearchByStr(ip)
	if err != nil {
		return result
	}

	split := strings.SplitN(r, "|", 5)
	result.Country = converEmpty(split[0])
	result.Region = converEmpty(split[2])
	result.City = converEmpty(split[3])
	result.ISP = converEmpty(split[4])

	return result
}

func converEmpty(str string) string {
	if str == "0" || str == "" {
		return "unknown"
	}
	return str
}
