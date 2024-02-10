package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/importing"
)

func NewImportFile() *cobra.Command {
	var o string
	var fps string
	var id string

	var cmdPrint = &cobra.Command{
		Use: "importfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			fpsAsInt, err := strconv.Atoi(fps)
			if err != nil {
				return err
			}

			// Create the lock file and db
			lockFile := fmt.Sprintf("%v/%v", o, "db.lock")
			db, err := db.OpenDb[importing.DbEntry](fmt.Sprintf("%v/%v", o, "db.json"), lockFile)
			if err != nil {
				return err
			}

			err = importing.ImportVideoFromPath(id, args[0], o, fpsAsInt, db)
			if err != nil {
				fmt.Println(err.(*errors.Error).ErrorStack())
				return err
			}
			return nil
		},
	}

	cmdPrint.Flags().StringVarP(&o, "output", "o", "", "Specify the output")
	cmdPrint.Flags().StringVar(&id, "id", "", "Specify the id")
	cmdPrint.MarkFlagRequired("id")
	cmdPrint.MarkFlagRequired("output")
	cmdPrint.MarkFlagRequired("fps")
	cmdPrint.Flags().StringVar(&fps, "fps", "", "Specify the wanted fps")

	return cmdPrint
}
