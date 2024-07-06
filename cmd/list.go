package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	pkg "github.com/cnnvito/dnsregion"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:    "list",
	Short:  "List subnet resource",
	PreRun: func(cmd *cobra.Command, args []string) { initSubnet() },
	Run: func(cmd *cobra.Command, args []string) {
		s := strings.Builder{}
		s.WriteString("当前包含运营商列表: \n")

		for _, g := range subnetClieng.Groups {
			s.WriteString("* ")
			s.WriteString(g.Name)
			for _, c := range subnetClieng.SearchGroupChildrenNodes(g.Name) {
				s.WriteString("\n|- ")
				s.WriteString(c.Name)
				s.WriteString("\t- ")
				s.WriteString(c.Ip)
			}
			s.WriteString("\n\n")
		}

		fmt.Println(s.String())
	},
}

var (
	subnetFile   string
	subnetClieng *pkg.SubnetResoucre
)

func initListCmd() {
	listCmd.Flags().StringVar(&subnetFile, "subnet-file", "", "Path to the subnet file")
}

func initSubnet() {
	if subnetFile != "" {
		fs, err := os.Open(subnetFile)
		if err != nil {
			panic(err)
		}

		r, err := io.ReadAll(fs)
		if err != nil {
			panic(err)
		}

		subnetClieng, err = pkg.NewSubnetResouce(r)
		if err != nil {
			panic(err)
		}
	} else {
		subnetClieng = pkg.DefaultSubnetResource
	}
}
