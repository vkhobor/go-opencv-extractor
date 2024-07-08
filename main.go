package main

import (
	"github.com/vkhobor/go-opencv/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "go-extractor"}

	rootCmd.AddCommand(cmd.NewConfig())
	rootCmd.AddCommand(cmd.NewVersion())
	rootCmd.AddCommand(cmd.NewRunserver())
	rootCmd.AddCommand(cmd.NewRunOpenApi())

	cobra.CheckErr(rootCmd.Execute())
}
