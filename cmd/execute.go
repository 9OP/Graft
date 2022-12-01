package cmd

import (
	"fmt"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

// var (
// 	consistency bool
// )

func execute(entry string, logType domain.LogType, cluster string) error {
	res, err := pkg.Execute(entry, logType, cluster)
	if err != nil {
		return err
	}
	if res.Err != nil {
		fmt.Println("err", res.Err)
		return res.Err
	}
	fmt.Println(string(res.Out))
	return nil
}

var queryCmd = &cobra.Command{
	Use:     "query [entry]",
	GroupID: "fsm",
	Short:   "Execute query (read-only)",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute(args[0], domain.LogQuery, cluster.String())
	},
}

var commandCmd = &cobra.Command{
	Use:     "command [entry]",
	GroupID: "fsm",
	Short:   "Execute command (write-only)",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute(args[0], domain.LogCommand, cluster.String())
	},
}

func init() {
	// queryCmd.Flags().BoolVarP(&consistency, "consistency", "c", true, "Require strong consistency. Only for QUERY.")
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(commandCmd)
}
