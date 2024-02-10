package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/importing"
)

func NewImport() *cobra.Command {
	var o string
	var fps string

	var cmdPrint = &cobra.Command{
		Use: "importurl",
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

			err = importing.ImportVideo(args[0], o, fpsAsInt, db)
			if err != nil {
				fmt.Println(err.(*errors.Error).ErrorStack())
				return err
			}
			return nil
		},
	}

	cmdPrint.Flags().StringVar(&o, "output", "", "Specify the output")
	cmdPrint.Flags().StringVar(&fps, "fps", "", "Specify the wanted fps")

	return cmdPrint
}
