package dnsregion

import "net"

type DnsClientInterface interface {
	GetServer() string
	Lookup(domain, qtype string) ([]string, error)
	LookupWithSubnet(domain, qtype, addr string) ([]string, error)
}

type DnsRegionResult struct {
	Info  []IPResult
	Error error
}

type queryRecordOption func(qr *DnsRegionInterface)

func WithDnsClientOption(dnsClient DnsClientInterface) queryRecordOption {
	return func(qr *DnsRegionInterface) {
		if dnsClient != nil {
			qr.DnsClient = dnsClient
		}
	}
}

func WithIPClientOption(ipClient *IPDatabase) queryRecordOption {
	return func(qr *DnsRegionInterface) {
		if ipClient != nil {
			qr.IPClient = ipClient
		}
	}
}

func NewQueryRecord(opts ...queryRecordOption) *DnsRegionInterface {
	qr := &DnsRegionInterface{}

	for _, opt := range opts {
		opt(qr)
	}

	if qr.DnsClient == nil {
		qr.DnsClient = NewDNSClient()
	}

	if qr.IPClient == nil {
		qr.IPClient = DefaultIPDatabase
	}

	return qr
}

type DnsRegionInterface struct {
	DnsClient DnsClientInterface
	IPClient  *IPDatabase
}

func (c *DnsRegionInterface) Query(domain, qtype string) DnsRegionResult {
	var result DnsRegionResult
	rst, err := c.DnsClient.Lookup(domain, qtype)
	if err != nil {
		result.Error = err
		return result
	}

	result.Info = c.processResult(rst)
	return result
}

func (c *DnsRegionInterface) QueryWithSubnet(domain, qtype, subnet string) DnsRegionResult {
	var result DnsRegionResult
	rst, err := c.DnsClient.LookupWithSubnet(domain, qtype, subnet)
	if err != nil {
		result.Error = err
		return result
	}

	result.Info = c.processResult(rst)
	return result
}

func (c *DnsRegionInterface) processResult(data []string) []IPResult {
	rst := make([]IPResult, 0, len(data)) // cname domain?
	for _, r := range data {
		ip := net.ParseIP(r)
		if ip == nil {
			continue
		}
		rst = append(rst, c.IPClient.Parser(r))
	}
	return rst
}
