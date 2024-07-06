package main

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/spf13/cobra"

	pkg "github.com/cnnvito/dnsregion"
)

var baseCmd = &cobra.Command{}

var (
	dnsServer      string
	dnsPort        int
	isTls          bool
	isDOH          bool
	ipDatabasefile string

	queryClient *pkg.DnsRegionInterface
)

func initBaseCmd() {
	baseCmd.Flags().StringVar(&dnsServer, "dns-server", "", "dns server")
	baseCmd.Flags().IntVar(&dnsPort, "dns-port", 0, "dns port")
	baseCmd.Flags().BoolVar(&isTls, "use-tls", false, "use tcp-tls")
	baseCmd.Flags().BoolVar(&isDOH, "use-doh", false, "use doh")
	baseCmd.Flags().StringVar(&ipDatabasefile, "ipdb", "", "ip database file")
}

func initClient() {
	var (
		ipdb      *pkg.IPDatabase
		dnsClient pkg.DnsClientInterface
	)

	if isDOH {
		dnsClient = pkg.NewDnsDOHClient(pkg.WithDOHServer(dnsServer), pkg.WithDOHPort(dnsPort))
	} else {
		if isTls {
			dnsClient = pkg.NewDNSClient(pkg.WithDNSServer(dnsServer), pkg.WithServerPort(dnsPort), pkg.WithDOT())
		} else {
			dnsClient = pkg.NewDNSClient(pkg.WithDNSServer(dnsServer), pkg.WithServerPort(dnsPort))
		}
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
