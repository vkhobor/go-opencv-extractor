package main

import (
	"github.com/vkhobor/go-opencv/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "go-opencv"}

	rootCmd.AddCommand(cmd.NewCompare())
	rootCmd.AddCommand(cmd.NewImport())
	rootCmd.AddCommand(cmd.NewImportFile())
	rootCmd.AddCommand(cmd.NewCleanSpace())
	rootCmd.AddCommand(cmd.NewRunserver())

	cobra.CheckErr(rootCmd.Execute())
}
