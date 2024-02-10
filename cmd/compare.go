package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vkhobor/go-opencv/image"
)

func NewCompare() *cobra.Command {
	var cmdPrint = &cobra.Command{
		Use: "compare",
		RunE: func(cmd *cobra.Command, args []string) error {
			similarity, err := image.CompareImagesPath(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", similarity)
			return nil
		},
	}

	return cmdPrint
}
