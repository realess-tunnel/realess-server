package create

import (
	"log"
	"github.com/spf13/cobra"
	"realess-server/internal/network"
	"realess-server/internal/tun"
	"realess-server/pkg/utils"
)


/*
setupRules configures the necessary network rules for the TUN interface to function properly.

It enables IP forwarding, sets up NAT masquerading, and configures firewall rules to allow traffic between the TUN interface and the internet.
*/
func setupRules(tunName string, cidr string) {
	// Get the default physical interface to use for NAT
	physicalInterface, err := utils.GetPhysicalInterface()
	if err != nil {
		log.Fatalf("Error getting physical interface: %v", err)
	}

	// Enable IP forwarding
	if err := network.EnableIPForwarding(); err != nil {
		log.Fatalf("Error enabling IP forwarding: %v", err)
	}

	// Enable NAT masquerading
	if err := network.SetupNatMasquerade(cidr, physicalInterface); err != nil {
		log.Fatalf("Error enabling NAT masquerading: %v", err)
	}

	// Configure firewall rules to allow traffic between TUN interface and the internet
	if err := network.SetupForwardRules(physicalInterface, tunName); err != nil {
		log.Fatalf("Error configuring firewall: %v", err)
	}
}


func rlssCreate(tunName string, cidr string) {
	// Setup network rules for the TUN interface
	setupRules(tunName, cidr)

	// Create the TUN interface
	serverTun := tun.CreateTunInterface(tunName, cidr)

	tun.StartTunServer(serverTun)

}

/*
CreateCreateCmd represents the create command to install the application and dependencies.
It sets up the initial configuration and runs the service.
*/
func CreateCreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:     "create [tun-name]",
		Aliases: []string{"install", "setup"},
		Short:   "Create a realess server instance with the specified tun interface name (default: realess-tun).",
		Args:    cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Setup the server configuration
			tunName := "realess-tun"
			if len(args) > 0 {
				tunName = args[0]
			}

			rlssCreate(tunName, "10.0.0.0/24")
		},
	}

	return createCmd
}
