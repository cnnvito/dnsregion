package dnsregion

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
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

func WithDOHPort(port int) dnsDOHClientOption {
	return func(c *dnsDOHClient) {
		if port > 0 {
			c.port = strconv.Itoa(port)
		} else {
			c.port = "443"
		}
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
	port   string

	o sync.Once
}

func (c *dnsDOHClient) GetServer() string {
	return c.serverAddr()
}

func (c *dnsDOHClient) ResolveValues(domain, qtype string) ([]string, error) {
	params := url.Values{}
	params.Set("name", domain)
	params.Set("type", qtype)
	return c.query(&params)
}

// Not all DoH supports edns_client_subnet parameter.
// for example, Cloudflare's DoH doesn't support it.
func (c *dnsDOHClient) ResolveValuesWithSubnet(domain, qtype, subnet string) ([]string, error) {
	params := url.Values{}
	params.Set("name", domain)
	params.Set("type", qtype)
	params.Set("edns_client_subnet", subnet)
	return c.query(&params)
}

func (c *dnsDOHClient) query(params *url.Values) (result []string, err error) {
	p := params.Encode()
	req, err := http.NewRequest("GET", c.serverAddr()+"?"+p, nil)
	if err != nil {
		return
	}
	req.Header.Set("accept", "application/dns-json")

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

	var data dohResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return result, err
	}

	for _, r := range data.Answer {
		result = append(result, r.Data)
	}
	return
}

func (c *dnsDOHClient) serverAddr() string {
	c.o.Do(func() {
		if c.server == "" {
			c.server = defaultDOHServer
			return
		}

		u, _ := url.Parse(c.server)
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		if c.port != "" && u.Port() == "" {
			u.Host = u.Host + ":" + c.port
		}
		c.server = defaultDOHServer
	})
	return c.server
}

type dohResponse struct {
	Status int  `json:"Status"`
	Tc     bool `json:"TC"`
	Rd     bool `json:"RD"`
	Ra     bool `json:"RA"`
	Ad     bool `json:"AD"`
	Cd     bool `json:"CD"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}
