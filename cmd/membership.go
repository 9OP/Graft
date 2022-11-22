package cmd

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"graft/pkg"
	"graft/pkg/domain"

	"github.com/spf13/cobra"
)

var (
	peer         domain.Peer
	addr         net.IP
	port         uint16
	updateConfig bool
)

var membershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Manage cluster memberships configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add ['id|host|p2p-port|api-port']",
	Short: "Add new node to cluster",
	Long: `
Add new node to cluster:

	Add new node to cluster configuration
	Peer is first added as inactive. When all the other nodes are aware of the
	configuration change, then node is marked as active.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		// format: <id>|<ip>|<p2p-port>|<api-port>
		re := regexp.MustCompile(`^(^[a-z-A-Z][a-z-A-Z-0-9]*)\|(.*)?\|[0-9]{1,5}\|[0-9]{1,5}$`)
		match := re.Find([]byte(args[0]))
		if match == nil {
			return fmt.Errorf("\n\node format unknown %v, should use format: <id>|<host>|<p2p-port>|<api-port>", args[0])
		}

		v := strings.Split(args[0], "|")

		newPeer, err := domain.NewPeer(
			v[0],
			false,
			v[1],
			v[2],
			v[3],
		)
		if err != nil {
			return err
		}

		peer = *newPeer

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		pkg.AddMember(peer, addr, port)
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove [node-id]",
	Short: "Remove old node from cluster",
	Long: `
Remove old node from cluster:

	Remove node-id from cluster configuration.
	Peer-id is first marked as inactive. When all the other nodes are aware of the
	configuration change, then node-id is removed from the configuration.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		pkg.RemoveMember(id, addr, port)
	},
}

func init() {
	membershipCmd.PersistentFlags().BoolVarP(&updateConfig, "change", "c", false, "Update configuration file")
	membershipCmd.PersistentFlags().IPVarP(&addr, "addr", "a", nil, "Host of a running cluster node")
	membershipCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 8080, "Port of a running cluster node")

	membershipCmd.MarkPersistentFlagRequired("addr")

	membershipCmd.AddCommand(addMemberCmd)
	membershipCmd.AddCommand(removeMemberCmd)
	rootCmd.AddCommand(membershipCmd)
}
