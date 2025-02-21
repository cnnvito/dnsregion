package main

import (
	"fmt"
	"os"

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

		fmt.Fprintf(os.Stdout, "DNS SERVER: %s\n", queryClient.DnsClient.GetServer())

		for _, domain := range args {
			var result pkg.DnsRegionResult

			fmt.Fprintf(os.Stdout, "\n* %s", domain)
			if subnet == "" {
				result = queryClient.Query(domain, qtype)
			} else {
				result = queryClient.QueryWithSubnet(domain, qtype, subnet)
				fmt.Fprintf(os.Stdout, " (edns_client_subnet: %s)", subnet)
			}
			fmt.Fprintf(os.Stdout, "\n")
			for _, r := range result.Info {
				fmt.Fprintf(os.Stdout, "|- %2s%25s\n", r.Ip, r.Region)
			}
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "|- ERROR: %s\n", result.Error.Error())
			}
		}
	},
}

func initDnsCmd() {
	dnsCmd.Flags().StringP("type", "t", "A", "query dns type")
	dnsCmd.Flags().StringP("subnet", "s", "", "use edns_client_subnet protocol")
	dnsCmd.Flags().AddFlagSet(baseCmd.Flags())
}
