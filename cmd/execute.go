package cmd

import (
	"fmt"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var (
	exType      executeType
	consistency bool
)

var executeCmd = &cobra.Command{
	Use:   "execute [entry]",
	Short: "Execute entry on FSM",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		entry := args[0]

		var logType domain.LogType
		switch exType {
		case "COMMAND":
			logType = domain.LogCommand
		case "QUERY":
			logType = domain.LogQuery
		}

		res, err := pkg.Execute(entry, logType, cluster.String())
		if err != nil {
			return err
		}

		if res.Err != nil {
			fmt.Println("err", res.Err)
			return res.Err
		}

		fmt.Println(string(res.Out))

		return nil
	},
}

func init() {
	executeCmd.Flags().BoolVarP(&consistency, "consistency", "c", true, "Require strong consistency. Only for QUERY.")
	executeCmd.Flags().Var(&exType, "type", `Execute type`)
	executeCmd.MarkFlagRequired("type")
	clusterCmd.AddCommand(executeCmd)
}
