package cmd

import (
	"fmt"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var exType executeType

var executeCmd = &cobra.Command{
	Use:   "execute [command]",
	Short: "Execute command on FSM",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterPeer := domain.Peer{Host: cluster.AddrPort}
		entry := args[0]

		var logType domain.LogType
		switch exType {
		case "COMMAND":
			logType = domain.LogCommand
		case "QUERY":
			logType = domain.LogQuery
		}

		res, err := pkg.Execute(entry, logType, clusterPeer)
		if err != nil {
			return err
		}

		if res.Err != nil {
			return fmt.Errorf("fsm eval error:\n\t%v", res.Err.Error())
		}

		fmt.Println(string(res.Out))

		return nil
	},
}

func init() {
	executeCmd.Flags().Var(&exType, "type", `Execute type`)
	executeCmd.MarkFlagRequired("type")

	clusterCmd.AddCommand(executeCmd)
}
