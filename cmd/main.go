package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnsregion",
	Short: "dnsregion is a dns region detect tool",
}

func init() {
	// init order
	initBaseCmd()
	initListCmd()
	initDnsCmd()
	initEdnsCmd()

	rootCmd.AddCommand(dnsCmd, listCmd, ednsCmd)
}

func main() {
	rootCmd.Execute()
}
