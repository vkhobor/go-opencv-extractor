package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/vkhobor/go-opencv/config"
)

func NewConfig() *cobra.Command {
	output := ""

	var cmdPrint = &cobra.Command{
		Use: "config",
		RunE: func(cmd *cobra.Command, args []string) error {

			viperconf := config.MustNewDefaultViperConfig()

			slog.Info("Writing config", "output", output)
			err := viperconf.WriteConfigAs(output)
			if err != nil {
				return err
			}

			// give some feedback to the user where to put the config files on the system
			fmt.Println("You can put config files in viper supported formats in the following locations")

			for _, path := range config.ConfigPaths {
				fmt.Println(path)
			}

			return nil
		},
	}

	cmdPrint.Flags().StringVarP(&output, "output", "o", "", "Specify the output file")
	cmdPrint.MarkFlagRequired("output")

	return cmdPrint
}
