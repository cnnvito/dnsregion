package main

import (
	"fmt"
	"strings"

	pkg "github.com/cnnvito/dnsregion"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:        "dns [domains...]",
	Short:      "Domain Name System (DNS) resolution",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"domain"},
	PreRun: func(cmd *cobra.Command, args []string) {
		initClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		qtype, _ := cmd.Flags().GetString("type")
		subnet, _ := cmd.Flags().GetString("subnet")

		s := strings.Builder{}

		s.WriteString("DNS SERVER: ")
		s.WriteString(queryClient.DnsClient.GetServer())
		s.WriteString("\n")

		for _, domain := range args {
			var result pkg.DnsRegionResult

			s.WriteString("\n")
			s.WriteString("* ")
			s.WriteString(domain)

			if subnet == "" {
				result = queryClient.Query(domain, qtype)
			} else {
				result = queryClient.QueryWithSubnet(domain, qtype, subnet)
				s.WriteString(" (edns_client_subnet: ")
				s.WriteString(subnet)
				s.WriteString(")")
			}

			s.WriteString("\n")
			for _, r := range result.Info {
				s.WriteString("|- ")
				s.WriteString(r.IP)
				s.WriteString("\t")
				s.WriteString(r.Country)
				s.WriteString("|")
				s.WriteString(r.Region)
				s.WriteString("|")
				s.WriteString(r.City)
				s.WriteString("|")
				s.WriteString(r.ISP)
				s.WriteString("\n")
			}
			if result.Error != nil {
				s.WriteString("|-")
				s.WriteString(result.Error.Error())
				s.WriteString("\n")
			}
		}
		fmt.Println(s.String())
	},
}

func initDnsCmd() {
	dnsCmd.Flags().StringP("type", "t", "A", "query dns type")
	dnsCmd.Flags().StringP("subnet", "s", "", "use edns_client_subnet protocol")
	dnsCmd.Flags().AddFlagSet(baseCmd.Flags())
}
