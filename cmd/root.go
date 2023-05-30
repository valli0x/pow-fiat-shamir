package cmd

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use:   "pow-fiat-shamir",
		Short: "multi-sign needed for send eth transaction with multi sign",
	}
)
