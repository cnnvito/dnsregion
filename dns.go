package dnsregion

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/miekg/dns"
)

const (
	defaultDnsServer = "223.5.5.5"
	defaultDnsPort   = "53"
	defaultDOTPort   = "853"
)

var QTypes = map[string]uint16{
	"A":     dns.TypeA,
	"AAAA":  dns.TypeAAAA,
	"CNAME": dns.TypeCNAME,
	"MX":    dns.TypeMX,
	"NS":    dns.TypeNS,
	"PTR":   dns.TypePTR,
}

type dnsClientOption func(c *DnsClient)

func WithDNSServer(server string) dnsClientOption {
	return func(c *DnsClient) {
		c.server = server
	}
}

func WithServerPort(port int) dnsClientOption {
	return func(c *DnsClient) {
		if port > 0 {
			c.port = strconv.Itoa(port)
		}
	}
}

func WithDOT() dnsClientOption {
	return func(c *DnsClient) {
		c.isTLS = true
		c.Net = "tcp-tls"
	}
}

func NewDNSClient(opts ...dnsClientOption) *DnsClient {
	c := &DnsClient{
		Client: &dns.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type DnsClient struct {
	*dns.Client

	server string
	port   string
	isTLS  bool

	o sync.Once
}

func (c *DnsClient) serverAddress() string {
	c.o.Do(func() {
		port := defaultDnsPort
		server := defaultDnsServer
		if c.isTLS && c.port == "" {
			port = defaultDOTPort
		}
		if c.server != "" {
			server = c.server
		}

		c.server = net.JoinHostPort(server, port)
	})
	return c.server
}

func (c *DnsClient) GetServer() string {
	return c.serverAddress()
}

func (c *DnsClient) ResolveValues(domain, qtype string) ([]string, error) {
	if domain[len(domain)-1] != '.' {
		domain += "."
	}

	if qt, ok := QTypes[qtype]; ok { // 检查类型存在性
		m := new(dns.Msg)
		m.SetQuestion(domain, qt)

		return c.query(m)
	} else {
		return nil, fmt.Errorf("invalid query type: %s", qtype)
	}
}

func (c *DnsClient) ResolveValuesWithSubnet(domain, qtype, addr string) ([]string, error) {
	if domain[len(domain)-1] != '.' {
		domain += "."
	}

	if _, ok := QTypes[qtype]; !ok { // 检查类型存在性
		return nil, fmt.Errorf("invalid query type: %s", qtype)
	}

	clientA := net.ParseIP(addr)
	if clientA == nil {
		return []string{}, fmt.Errorf("invalid IP address: %s", addr)
	}

	var is4 uint16 = 1
	var is4mask uint8 = 24
	if clientA.To4() == nil {
		is4 = 2
		is4mask = 128
	}

	opt := &dns.OPT{
		Hdr: dns.RR_Header{
			Name:   ".",
			Rrtype: dns.TypeOPT,
		},
		Option: []dns.EDNS0{
			&dns.EDNS0_SUBNET{
				Code:          dns.EDNS0SUBNET,
				Family:        is4,
				SourceNetmask: is4mask,
				SourceScope:   0,
				Address:       clientA,
			},
		},
	}

	m := new(dns.Msg)
	m.Extra = append(m.Extra, opt)
	m.SetQuestion(domain, QTypes[qtype])

	return c.query(m)
}

func (c *DnsClient) query(msg *dns.Msg) ([]string, error) {
	r, _, err := c.Exchange(msg, c.serverAddress())
	if err != nil {
		return []string{}, err
	}
	return parseValues(r.Answer), nil
}

func parseValues(rr []dns.RR) []string {
	lenRR := len(rr)
	if lenRR == 0 {
		return []string{}
	}

	results := make([]string, 0, lenRR)
	for _, r := range rr {
		h := r.Header()
		switch h.Rrtype {
		case dns.TypeA:
			a := r.(*dns.A)
			results = append(results, a.A.String())
		case dns.TypeAAAA:
			a := r.(*dns.AAAA)
			results = append(results, a.AAAA.String())
		case dns.TypeCNAME:
			c := r.(*dns.CNAME)
			results = append(results, c.Target)
		case dns.TypeMX:
			m := r.(*dns.MX)
			results = append(results, m.Mx)
		case dns.TypeNS:
			n := r.(*dns.NS)
			results = append(results, n.Ns)
		case dns.TypePTR:
			p := r.(*dns.PTR)
			results = append(results, p.Ptr)
		}
	}
	return results
}
