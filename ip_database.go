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
	Ip     string
	Region string
}

func (i IPResult) String() string {
	r, _ := json.Marshal(&i)
	return string(r)
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
	r, err := s.SearchByStr(ip)
	if err != nil {
		return IPResult{}
	}
	rs := []string{}
	for _, el := range strings.Split(r, "|") {
		rs = append(rs, converEmpty(el))
	}
	return IPResult{Ip: ip, Region: strings.Join(rs, "|")}
}

func converEmpty(str string) string {
	if str == "0" || str == "" {
		return "N/A"
	}
	return str
}
