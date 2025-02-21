package dnsregion

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

const defaultDOHServer = "https://dns.alidns.com/resolve"

type dnsDOHClientOption func(c *dnsDOHClient)

func WithDOHServer(server string) dnsDOHClientOption {
	return func(c *dnsDOHClient) {
		if server == "" {
			server = defaultDOHServer
			return
		}
		c.server = server
	}
}

func NewDnsDOHClient(opts ...dnsDOHClientOption) *dnsDOHClient {
	c := &dnsDOHClient{
		client: &http.Client{
			Timeout: 15 * time.Second, // slow?
			Transport: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.server == "" {
		c.server = defaultDOHServer
	}
	return c
}

type dnsDOHClient struct {
	client *http.Client
	server string
}

func (c *dnsDOHClient) GetServer() string {
	return c.serverAddr()
}

func (c *dnsDOHClient) Lookup(domain, qtype string) ([]string, error) {
	m, err := MakeDnsQuestionMsg(domain, qtype, "")
	if err != nil {
		return nil, err
	}
	return c.query(m)
}

// Not all DoH supports edns_client_subnet parameter.
// for example, Cloudflare's DoH doesn't support it.
func (c *dnsDOHClient) LookupWithSubnet(domain, qtype, subnet string) ([]string, error) {
	m, err := MakeDnsQuestionMsg(domain, qtype, subnet)
	if err != nil {
		return nil, err
	}
	return c.query(m)
}

func (c *dnsDOHClient) query(dnsMsg *dns.Msg) (result []string, err error) {
	dnsPack, err := dnsMsg.Pack()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", c.serverAddr(), bytes.NewBuffer(dnsPack))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "dns-region/1.0")
	req.Header.Set("accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("dns query failed with status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	m := new(dns.Msg)
	if err := m.Unpack(body); err != nil {
		return nil, err
	}
	return parseValues(m.Answer), nil
}

func (c *dnsDOHClient) serverAddr() string {
	return c.server
}
