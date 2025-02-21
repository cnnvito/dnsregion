package main

import (
	"net"
	"strconv"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/spf13/cobra"

	pkg "github.com/cnnvito/dnsregion"
)

var baseCmd = &cobra.Command{}

var (
	dnsServer      string
	ipDatabasefile string

	queryClient *pkg.DnsRegionInterface
)

func initBaseCmd() {
	baseCmd.Flags().StringVar(&dnsServer, "dns-server", "", "dns server, like https://doh.pub/dns-query tls://223.5.5.5:853 tcp://223.5.5.5:53")
	baseCmd.Flags().StringVar(&ipDatabasefile, "ipdb", "", "ip database file")
}

func initClient() {
	var (
		ipdb      *pkg.IPDatabase
		dnsClient pkg.DnsClientInterface
	)

	scheme, host, _port := splitDnsServer(dnsServer)
	port, _ := strconv.Atoi(_port)
	switch scheme {
	case "http", "https":
		dnsClient = pkg.NewDnsDOHClient(pkg.WithDOHServer(dnsServer))
	case "tls":
		dnsClient = pkg.NewDNSClient(pkg.WithDNSServer(host), pkg.WithServerPort(port), pkg.WithDOT())
	case "tcp":
		dnsClient = pkg.NewDNSClient(pkg.WithDNSServer(host), pkg.WithServerPort(port), pkg.WithTcp())
	default:
		dnsClient = pkg.NewDNSClient(pkg.WithDNSServer(host), pkg.WithServerPort(port))
	}

	if ipDatabasefile != "" {
		buf, err := xdb.LoadContentFromFile(ipDatabasefile)
		if err != nil {
			panic(err)
		}
		ipdb = pkg.NewIPDatabase(buf)
	} else {
		ipdb = pkg.DefaultIPDatabase
	}

	queryClient = pkg.NewQueryRecord(pkg.WithDnsClientOption(dnsClient), pkg.WithIPClientOption(ipdb))
}

func splitDnsServer(addr string) (scheme, host, port string) {
	s := strings.Split(addr, "://")
	if len(s) > 1 {
		host = s[1]
		scheme = s[0]
	} else {
		host = s[0]
	}

	_host, _port, err := net.SplitHostPort(host)
	if err == nil {
		host = _host
		port = _port
	}
	return
}
