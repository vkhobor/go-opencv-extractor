package cmd

import (
	"database/sql"

	"github.com/spf13/cobra"
)

func NewImportVideoFromFolder() *cobra.Command {
	var o string

	var cmdPrint = &cobra.Command{
		Use: "import-folder-video",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := sql.Open("sqlite3", o)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmdPrint.Flags().StringVarP(&o, "output-sqlite", "o", "", "Specify the output")
	cmdPrint.MarkFlagRequired("output")

	return cmdPrint
}
