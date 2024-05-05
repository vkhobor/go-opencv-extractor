package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersion() *cobra.Command {

	var cmdPrint = &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Version 1.0.0")
			return nil
		},
	}

	return cmdPrint
}
