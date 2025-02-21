package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ednsCmd = &cobra.Command{
	Use:        "edns [group...]",
	Short:      "batch client subnet query record with group",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"group"},
	PreRun: func(cmd *cobra.Command, args []string) {
		initClient()
		initSubnet()
	},
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		qtype, _ := cmd.Flags().GetString("type")

		fmt.Fprintf(os.Stdout, "Dns Server:\t%s\nQuery:\t%s %s\n\n", queryClient.DnsClient.GetServer(), qtype, domain)
		for _, name := range args {
			fmt.Fprintf(os.Stdout, "* %s\n", name)

			for _, node := range subnetClieng.SearchGroupChildrenNodes(name) {
				fmt.Fprintf(os.Stdout, "|- * %s (edns_client_subnet=%s)\n", node.Name, node.Ip)

				tempRst := queryClient.QueryWithSubnet(domain, qtype, node.Ip)
				for _, ip := range tempRst.Info {
					fmt.Fprintf(os.Stdout, "|  |- %s(%s)\n", ip.Region, ip.Ip)
				}

				if tempRst.Error != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", tempRst.Error.Error())
				}
			}
		}
	},
}

func initEdnsCmd() {
	ednsCmd.Flags().StringP("domain", "d", "", "query domain name")
	ednsCmd.Flags().StringP("type", "t", "A", "query domain type")

	ednsCmd.Flags().AddFlagSet(baseCmd.Flags())
	ednsCmd.Flags().AddFlagSet(listCmd.Flags())

	ednsCmd.MarkFlagRequired("domain")
}
