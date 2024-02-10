package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/importing"
)

func NewCleanSpace() *cobra.Command {
	var o string

	var cmdPrint = &cobra.Command{
		Use: "cleanspace",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Create the lock file and db
			lockFile := fmt.Sprintf("%v/%v", o, "db.lock")
			db, err := db.OpenDb[importing.DbEntry](fmt.Sprintf("%v/%v", o, "db.json"), lockFile)
			if err != nil {
				return err
			}

			all, err := db.ReadAll(context.Background())

			if err != nil {
				return err
			}

			allFilenames := make([]string, 0)
			for _, value := range all {
				if value.Status == importing.StatusImported {
					allFilenames = append(allFilenames, value.FileNames...)
				}
			}

			fmt.Println(len(allFilenames))

			entries, err := os.ReadDir(o)
			if err != nil {
				return err
			}
			filtered := lo.Filter(entries, func(entry fs.DirEntry, _ int) bool {
				return !(entry.Name() == "db.json" || entry.Name() == "db.lock")
			})

			for _, entry := range filtered {
				if !contains(allFilenames, entry.Name()) {
					err := os.Remove(fmt.Sprintf("%v/%v", o, entry.Name()))
					if err != nil {
						return err
					}
				}
			}

			for key, value := range all {
				if value.Status == importing.StatusError {
					err := db.Put(context.Background(), key, importing.DbEntry{
						Status: importing.StatusZero,
						Title:  value.Title,
						Url:    value.Url,
					})
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	cmdPrint.Flags().StringVarP(&o, "output", "o", "", "Specify the output")
	cmdPrint.MarkFlagRequired("output")

	return cmdPrint
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
