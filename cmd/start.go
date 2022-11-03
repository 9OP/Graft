package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
Options for start:
- ARG: which id of the cluster to start
- FLAG config path
- FLAG log level

*/

var configPath string
var logLevel string

var startCmd = &cobra.Command{
	Use:   "start [peer-id]",
	Short: "Start a cluster node",
	Args:  cobra.ExactArgs(1),
	// Should error when id is not in the config
	// By default will look for {pwd}/graft-config.yaml
	//
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	// Optionally run one of the validators provided by cobra
	// 	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
	// 		return err
	// 	}
	// 	// Run the custom validation logic
	// 	if myapp.IsValidColor(args[0]) {
	// 	  return nil
	// 	}
	// 	return fmt.Errorf("invalid color specified: %s", args[0])
	//   },
	Run: func(cmd *cobra.Command, args []string) {
		// res := stringer.Reverse(args[0])
		// fmt.Println(res)
		fmt.Println("start cluster node")
		confFlag, _ := cmd.Flags().GetString("config")
		fmt.Println(confFlag)
	},
}

func init() {
	// For log level see:
	// https://stackoverflow.com/questions/50824554/permitted-flag-values-for-cobra
	startCmd.Flags().StringVarP(&configPath, "config", "c", "graft-config.yaml", "Configuration file path")
	rootCmd.AddCommand(startCmd)
}

// var onlyDigits bool
// var inspectCmd = &cobra.Command{
//     Use:   "inspect",
//     Aliases: []string{"insp"},
//     Short:  "Inspects a string",
//     Args:  cobra.ExactArgs(1),
//     Run: func(cmd *cobra.Command, args []string) {

//         i := args[0]
//         res, kind := stringer.Inspect(i, onlyDigits)

//         pluralS := "s"
//         if res == 1 {
//             pluralS = ""
//         }
//         fmt.Printf("'%s' has %d %s%s.\n", i, res, kind, pluralS)
//     },
// }

// func init() {
//     inspectCmd.Flags().BoolVarP(&onlyDigits, "digits", "d", false, "Count only digits")
//     rootCmd.AddCommand(inspectCmd)
// }
