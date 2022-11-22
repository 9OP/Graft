package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var peer domain.Peer

var membershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Manage cluster memberships configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add ['id|port|p2p-port|api-port']",
	Short: "Add new node to cluster",
	Long: `
Add new node to cluster:

	Add new node to cluster configuration
	Peer is first added as inactive. When all the other nodes are aware of the
	configuration change, then peer is marked as active.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		// format: <id>|<ip>|<p2p-port>|<api-port>
		re := regexp.MustCompile(`^(^[a-z-A-Z][a-z-A-Z-0-9]*)\|([0-9]{1,3}.){3}.([0-9]{1,3})\|[0-9]{1,5}\|[0-9]{1,5}$`)
		match := strings.Split(string(re.Find([]byte(args[0]))), "|")
		if len(match) == 0 {
			return fmt.Errorf("\n\tpeer format unknown %v, should use format: <id>|<ip>|<p2p-port>|<api-port>", args[0])
		}

		newPeer, err := domain.NewPeer(
			match[0],
			false,
			match[1],
			match[2],
			match[3],
		)
		if err != nil {
			return err
		}

		peer = *newPeer

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add node", peer)
		// pkg.addMember(peer)
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove [peer-id]",
	Short: "Remove old node from cluster",
	Long: `
Remove old node from cluster:

	Remove peer-id from cluster configuration.
	Peer-id is first marked as inactive. When all the other nodes are aware of the
	configuration change, then peer-id is removed from the configuration.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		fmt.Println("remove node", id)
		// pkg.removeMember(id)
	},
}

func init() {
	// startCmd.Flags().BoolVarP(&updateConfig, "change", "c", false, "Update configuration file")
	membershipCmd.AddCommand(addMemberCmd)
	membershipCmd.AddCommand(removeMemberCmd)
	rootCmd.AddCommand(membershipCmd)
}
