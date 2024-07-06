package main

import (
	"fmt"
	"strings"

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

		s := strings.Builder{}
		s.WriteString("Dns Server:\t")
		s.WriteString(queryClient.DnsClient.GetServer())
		s.WriteString("\n")
		s.WriteString("Query:\t")
		s.WriteString(qtype)
		s.WriteString(" ")
		s.WriteString(domain)
		s.WriteString("\n\n")
		for _, name := range args {
			s.WriteString("* ")
			s.WriteString(name)
			s.WriteString("\n")

			for _, node := range subnetClieng.SearchGroupChildrenNodes(name) {
				s.WriteString("|- * ")
				s.WriteString(node.Name)
				s.WriteString(" (edns_client_subnet=")
				s.WriteString(node.Ip)
				s.WriteString(")")
				s.WriteString("\n")

				tempRst := queryClient.QueryWithSubnet(domain, qtype, node.Ip)
				for _, ip := range tempRst.Info {
					s.WriteString("|  |- ")
					s.WriteString(ip.Country)
					s.WriteString("|")
					s.WriteString(ip.Region)
					s.WriteString("|")
					s.WriteString(ip.City)
					s.WriteString("(")
					s.WriteString(ip.IP)
					s.WriteString(")")
					s.WriteString("\n")
				}

				if tempRst.Error != nil {
					s.WriteString("   |- ")
					s.WriteString(tempRst.Error.Error())
					s.WriteString("\n")
				}
			}
			s.WriteString("\n")
		}
		fmt.Println(s.String())
	},
}

func initEdnsCmd() {
	ednsCmd.Flags().StringP("domain", "d", "", "query domain name")
	ednsCmd.Flags().StringP("type", "t", "A", "query domain type")

	ednsCmd.Flags().AddFlagSet(baseCmd.Flags())
	ednsCmd.Flags().AddFlagSet(listCmd.Flags())

	ednsCmd.MarkFlagRequired("domain")
}
