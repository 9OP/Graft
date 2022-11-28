package cmd

import (
	"fmt"

	"graft/pkg"

	"github.com/spf13/cobra"
)

var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "Print cluster configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := pkg.ClusterConfiguration(cluster.String())
		if err != nil {
			return fmt.Errorf("failed load configuration\n%v", err.Error())
		}

		c, err := configuration.ToJSON()
		if err != nil {
			return err
		}

		fmt.Println(string(c))

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(configurationCmd)
}
